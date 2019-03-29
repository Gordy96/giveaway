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
	"giveaway/instagram"
	"giveaway/instagram/account"
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

func (c *Client) GetUserInfo(id string) (*data.Owner, error) {
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

	dat, err := c.procResp(resp)

	if user, ok := dat["user"]; ok {
		u, ok := user.(map[string]interface{})
		if ok {
			return &data.Owner{Id: strconv.FormatInt(int64(u["pk"].(float64)), 10), Username: u["username"].(string)}, nil
		}
	}

	return nil, err
}

func (c *Client) IsFollower(o *data.Owner, id string) (bool, error) {
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
	c.acc.ResetCookies()
	c.setCookieJar(c.acc.MakeCookieJar())

	c.QeSync()
	c.LauncherSync()

	jsonBody := map[string]interface{}{
		"country_codes":       "[{\"country_code\":\"1\",\"source\":[\"me_profile\"]},{\"country_code\":\"1\",\"source\":[\"default\"]}]",
		"phone_id":            c.acc.PhoneId,
		"username":            c.acc.Username,
		"adid":                c.acc.AdId,
		"guid":                c.acc.GUID,
		"device_id":           c.acc.DeviceId,
		"google_tokens":       "",
		"password":            c.acc.Password,
		"login_attempt_count": "1",
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
	}

	return false, err
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
