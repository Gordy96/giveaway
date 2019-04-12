package api

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"giveaway/data"
	"giveaway/data/errors"
	"giveaway/data/owner"
	"giveaway/instagram"
	"giveaway/instagram/account"
	"giveaway/instagram/structures"
	"giveaway/instagram/structures/stories"
	"giveaway/utils"
	"hash"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
)

type Client struct {
	http     *http.Client
	acc      *account.Account
	uag      utils.IUserAgentGenerator
	loggedIn bool
}

func (c *Client) Sign(body []byte) string {
	key := instagram.Constants[instagram.Version].Key
	h := hmac.New(func() hash.Hash {
		return sha256.New()
	}, []byte(key))
	h.Write(body)
	return hex.EncodeToString(h.Sum(nil))
}

func (c *Client) newRequest(method string, url string, body io.Reader) *http.Request {
	r, _ := http.NewRequest(method, url, body)
	r.Header.Set("User-Agent", c.uag.Get())
	r.Header.Set("X-DEVICE-ID", c.acc.GUID)
	r.Header.Set("X-IG-App-ID", instagram.Constants[instagram.Version].AppID)
	r.Header.Set("X-IG-Connection-Speed", "-1kbps")
	r.Header.Set("X-IG-Bandwidth-Speed-KBPS", "-1.000")
	r.Header.Set("X-IG-Bandwidth-TotalBytes-B", "0")
	r.Header.Set("X-IG-Bandwidth-TotalTime-MS", "0")
	r.Header.Set("X-IG-Connection-Type", "WIFI")
	r.Header.Set("X-IG-Capabilities", "3brTvw==")
	r.Header.Set("X-FB-HTTP-Engine", "Liger")
	r.Header.Set("Accept-Language", "en-US")
	return r
}

func (c *Client) makeRequestString(jsonBody map[string]interface{}) string {

	q, _ := json.Marshal(jsonBody)
	query := string(c.Sign(q)) + "." + string(q)
	post := url.Values{}
	post.Set("signed_body", query)
	post.Set("ig_sig_key_version", "4")
	return post.Encode()
}

func (c *Client) findCSRF(uri *url.URL) string {
	var csrf = ""
	for _, cookie := range c.http.Jar.Cookies(uri) {
		if strings.Contains(cookie.Name, "csrf") {
			csrf = cookie.Value
			break
		}
	}
	return csrf
}

func (c *Client) procResp(resp *http.Response) (map[string]interface{}, error) {
	str, _ := ioutil.ReadAll(resp.Body)
	temp := map[string]interface{}{}
	var err error = nil

	err = json.Unmarshal([]byte(str), &temp)

	if err != nil {
		return temp, err
	}

	c.acc.SetCookies(resp.Request.URL, resp.Cookies())

	if resp.StatusCode != 200 {
		err = fmt.Errorf("%d, %s", resp.StatusCode, resp.Status)
	}
	if temp["status"].(string) != "ok" {
		err = fmt.Errorf("request failed")
	}
	return temp, err
}

func handleResponseErrors(response *http.Response) error {
	if response.StatusCode != 200 {
		switch response.StatusCode {
		case 403:
			return errors.HttpForbidden{}
		case 429:
			return errors.HttpTooManyRequests{}
		default:
			return fmt.Errorf("%d %s", response.StatusCode, response.Status)
		}
	}
	return nil
}

func (c *Client) GetUserInfo(id string) (*data.User, error) {
	uri := fmt.Sprintf("%s/api/v1/users/%s/info", instagram.AppHost, id)
	req := c.newRequest(
		"GET",
		uri,
		nil,
	)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	if err = handleResponseErrors(resp); err != nil {
		return nil, err
	}

	dat, err := c.procResp(resp)

	if user, ok := dat["user"]; ok {
		usr, ok := user.(map[string]interface{})
		if ok {
			u := &data.User{}
			u.Username = usr["username"].(string)
			u.Id = strconv.FormatInt(int64(usr["pk"].(float64)), 10)
			u.Follows = int64(usr["following_count"].(float64))
			u.Followers = int64(usr["follower_count"].(float64))
			u.IsBusiness = usr["is_business"].(bool)
			u.IsPrivate = usr["is_private"].(bool)
			u.IsVerified = usr["is_verified"].(bool)
			return u, nil
		}
	}

	return nil, err
}

func (c *Client) IsFollower(o *owner.Owner, id string) (bool, error) {
	uri := fmt.Sprintf("%s/api/v1/friendships/%s/followers/", instagram.AppHost, id)
	req := c.newRequest(
		"GET",
		uri,
		nil,
	)
	var err error = nil
	if o.Username == "" && o.Id != "" {
		so, err := c.GetUserInfo(o.Id)
		if err != nil {
			return false, err
		}
		o.Username = so.Username
	} else if o.Username == "" {
		return false, fmt.Errorf("owner object malformed")
	}

	query := req.URL.Query()
	query.Add("query", o.Username)
	query.Add("rank_token", uuid.New().String())
	req.URL.RawQuery = query.Encode()

	resp, err := c.http.Do(req)
	if err != nil {
		return false, err
	}

	if err = handleResponseErrors(resp); err != nil {
		return false, err
	}

	dat, err := c.procResp(resp)

	if m, ok := dat["message"]; ok && m.(string) == "login_required" {
		return false, errors.LoginRequired{}
	}

	var presence = false
	err = nil
	oid64, _ := strconv.ParseInt(o.Id, 10, 0)
	if users, ok := dat["users"]; ok {
		for _, user := range users.([]interface{}) {
			if int64(user.(map[string]interface{})["pk"].(float64)) == oid64 {
				presence = true
			}
		}
	}

	return presence, err
}

func (c *Client) QeSync() error {
	jsonBody := map[string]interface{}{
		"experiments": instagram.Constants[instagram.Version].Experiments,
	}
	if !c.loggedIn {
		jsonBody["id"] = c.acc.GUID
	} else {
		jsonBody["id"] = c.acc.Id
		jsonBody["_uid"] = c.acc.Id
		jsonBody["_uuid"] = c.acc.GUID
	}
	uri := instagram.AppHost + "/api/v1/qe/sync/"
	req := c.newRequest(
		"POST",
		uri,
		bytes.NewReader([]byte(c.makeRequestString(jsonBody))),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}

	_, err = c.procResp(resp)

	return err
}

func (c *Client) LauncherSync() error {
	jsonBody := map[string]interface{}{
		"configs": instagram.Constants[instagram.Version].Configs,
	}
	if !c.loggedIn {
		jsonBody["id"] = c.acc.GUID
	} else {
		jsonBody["id"] = c.acc.Id
		jsonBody["_uid"] = c.acc.Id
		jsonBody["_uuid"] = c.acc.GUID
	}

	uri, _ := url.Parse(instagram.AppHost + "/api/v1/launcher/sync/")

	csrf := c.findCSRF(uri)

	if csrf != "" {
		jsonBody["_csrftoken"] = csrf
	}

	req := c.newRequest(
		"POST",
		uri.String(),
		bytes.NewReader([]byte(c.makeRequestString(jsonBody))),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}

	_, err = c.procResp(resp)

	return err
}

func (c *Client) setCookieJar(jar http.CookieJar) {
	c.http.Jar = jar
}

func (c *Client) Login() (bool, error) {
	if c.acc == nil {
		return false, fmt.Errorf("account is nil")
	}
	c.acc.ResetCookies()
	c.setCookieJar(c.acc.MakeCookieJar())

	c.QeSync()
	c.LauncherSync()

	jsonBody := map[string]interface{}{
		"country_codes":       "[{\"country_code\":\"380\",\"source\":[\"me_profile\"]},{\"country_code\":\"44\",\"source\":[\"default\"]}]",
		"phone_id":            c.acc.PhoneId,
		"username":            c.acc.Username,
		"adid":                c.acc.AdId,
		"guid":                c.acc.GUID,
		"device_id":           c.acc.DeviceId,
		"google_tokens":       "[]",
		"password":            c.acc.Password,
		"login_attempt_count": "0",
	}

	uri, _ := url.Parse(instagram.AppHost + "/api/v1/accounts/login/")

	csrf := c.findCSRF(uri)

	if csrf != "" {
		jsonBody["_csrftoken"] = csrf
	}

	req := c.newRequest(
		"POST",
		uri.String(),
		bytes.NewReader([]byte(c.makeRequestString(jsonBody))),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	resp, err := c.http.Do(req)
	if err != nil {
		return false, err
	}

	temp, err := c.procResp(resp)

	if u, b := temp["logged_in_user"]; b {
		if c.acc.Id == "" {
			c.acc.Id = strconv.FormatInt(int64(u.(map[string]interface{})["pk"].(float64)), 10)
		}
		c.loggedIn = true
		return true, nil
	} else {
		b, _ := json.Marshal(temp)
		return false, fmt.Errorf("%s", string(b))
	}

	return false, err
}

func (c *Client) QueryHashTagStories(tag string, cb func(item stories.StoryItem) (bool, error)) (*structures.Story, error) {
	var story *structures.Story = nil
	r, err := c.GetUnseenHashTagStories(tag)
	if err != nil {
		return story, err
	}
	story = r.Story
	if story != nil {
		for _, i := range story.Items {
			if f, err := cb(i); !f {
				return story, err
			}
		}
	}
	return story, nil
}

func (c *Client) GetUnseenHashTagStories(hashTag string) (*structures.HashTagStoriesResponse, error) {
	var err error = nil
	uri := fmt.Sprintf("%s/api/v1/tags/%s/story/", instagram.AppHost, hashTag)
	req := c.newRequest(
		"GET",
		uri,
		nil,
	)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	if err = handleResponseErrors(resp); err != nil {
		return nil, err
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var dat = &structures.HashTagStoriesResponse{}
	err = json.Unmarshal(respBytes, dat)
	if err != nil {
		return nil, err
	}
	return dat, nil
}

func (c *Client) MarkHashTagStoriesAsSeen(belongs *structures.Story, stories ...structures.WatchedStoryEntry) (bool, error) {
	var err error = nil
	uri, err := url.Parse(fmt.Sprintf("%s/api/v2/media/seen/?reel=1&live_vod=0", instagram.AppHost))
	if err != nil {
		return false, err
	}

	jsonBody := map[string]interface{}{
		"_uid":               c.acc.Id,
		"_uuid":              c.acc.GUID,
		"container_module":   "hashtag_feed",
		"live_vods_skipped":  map[string]interface{}{},
		"nuxes_skipped":      map[string]interface{}{},
		"nuxes":              map[string]interface{}{},
		"reels":              structures.MakeWatchedStoryRequest(belongs, stories...),
		"reel_media_skipped": map[string]interface{}{},
	}

	csrf := c.findCSRF(uri)

	if csrf != "" {
		jsonBody["_csrftoken"] = csrf
	}

	req := c.newRequest(
		"POST",
		uri.String(),
		bytes.NewReader([]byte(c.makeRequestString(jsonBody))),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	resp, err := c.http.Do(req)
	if err != nil {
		return false, err
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var dat = map[string]interface{}{}
	err = json.Unmarshal(respBytes, &dat)
	if err != nil {
		return false, err
	}
	if status, ok := dat["status"]; ok && status.(string) == "ok" {
		return true, nil
	}
	return false, nil
}

func makeProxy(proxy string) *http.Transport {
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		log.Println(err)
	}

	return &http.Transport{Proxy: http.ProxyURL(proxyURL)}
}

func NewApiClient() *Client {
	cl := &Client{}
	cl.uag = AppUserAgentGenerator{}
	jar, _ := cookiejar.New(nil)

	httpClient := &http.Client{Jar: jar}

	cl.http = httpClient
	return cl
}

func NewProxiedApiClient(proxy string) *Client {
	cl := NewApiClient()
	cl.http.Transport = makeProxy(proxy)
	return cl
}

func NewApiClientWithAccount(acc *account.Account) *Client {
	cl := NewProxiedApiClient(acc.Proxy)
	cl.SetAccount(acc)
	return cl
}

func (c *Client) SetAccount(acc *account.Account) {
	c.acc = acc
	c.http.Transport = makeProxy(acc.Proxy)
	c.setCookieJar(c.acc.MakeCookieJar())
}

func (c *Client) GetAccount() *account.Account {
	return c.acc
}
