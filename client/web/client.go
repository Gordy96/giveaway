package web

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"giveaway/data"
	httpErrors "giveaway/data/errors"
	"giveaway/data/owner"
	"giveaway/instagram/account"
	"giveaway/instagram/structures"
	"giveaway/utils"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

func NewWebClient(generator utils.IUserAgentGenerator, proxy string) *Client {
	cl := &Client{uag: generator}
	jar, _ := cookiejar.New(nil)
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		log.Println(err)
	}
	httpClient := &http.Client{Jar: jar, Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}

	cl.http = httpClient
	return cl
}

type Client struct {
	http        *http.Client
	sig         string
	appId       string
	rolloutHash string
	acc         *account.Account
	uag         utils.IUserAgentGenerator
}

func (c *Client) getHttpClient() *http.Client {
	return c.http
}

func (c *Client) sign(r *http.Request) string {
	var strConcat string
	if v := r.URL.Query().Get("variables"); v != "" {
		strConcat = fmt.Sprintf("%s:%s", c.sig, v)
	} else {
		strConcat = fmt.Sprintf("%s:%s", c.sig, r.URL.Path)
	}
	raw := md5.Sum([]byte(strConcat))
	return hex.EncodeToString(raw[:])
}

func (c *Client) prepare(r *http.Request) {
	r.Header.Set("X-Instagram-GIS", c.sign(r))
	r.Header.Set("X-IG-App-ID", c.appId)
	r.Header.Set("X-Requested-With", "XMLHttpRequest")
}

func getResponseString(r *http.Response) (string, error) {
	respBytes, err := ioutil.ReadAll(r.Body)
	return string(respBytes), err
}

func (c *Client) makeRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.uag.Get())

	return req, nil
}

func (c *Client) Init() error {
	req, err := c.makeRequest("GET", "https://www.instagram.com", nil)
	if err != nil {
		return err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	respString, err := getResponseString(resp)
	if err != nil {
		return err
	}

	sigStr := respString[strings.Index(respString, "rhx_gis\":\"")+10:]
	sigStr = sigStr[:strings.Index(sigStr, "\"")]

	c.sig = sigStr

	roll := respString[strings.Index(respString, "rollout_hash\":\"")+15:]
	roll = roll[:strings.Index(roll, "\"")]

	c.rolloutHash = roll

	consumerCommonsHashPositionStart := strings.Index(respString, "ConsumerCommons.js/") + 19
	respString = respString[consumerCommonsHashPositionStart:]
	consumerCommonsHashPositionEnd := strings.Index(respString, "\"")
	consumerCommonsHref := "/static/bundles/metro/ConsumerCommons.js/" + respString[:consumerCommonsHashPositionEnd]

	nUrl, _ := req.URL.Parse(consumerCommonsHref)
	req, err = c.makeRequest("GET", nUrl.String(), nil)

	if err != nil {
		return err
	}

	resp, err = c.http.Do(req)
	if err != nil {
		return nil
	}
	respString, err = getResponseString(resp)
	if err != nil {
		return nil
	}
	respString = respString[strings.Index(respString, "instagramWebFBAppId='")+21:]
	appId := respString[:strings.Index(respString, "'")]
	c.appId = appId

	return nil
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

func (c *Client) Login() (bool, error) {
	post := url.Values{}
	post.Set("username", c.acc.Username)
	post.Set("queryParams", "{}")
	post.Set("password", c.acc.Password)
	post.Set("optIntoOneTap", "false")
	uri, _ := url.Parse("https://www.instagram.com/accounts/login/ajax/")
	req, err := c.makeRequest("POST", uri.String(), bytes.NewReader([]byte(post.Encode())))
	if err != nil {
		return false, err
	}
	c.prepare(req)
	req.Header.Del("X-Instagram-GIS")
	req.Header.Set("Referer", "https://www.instagram.com/accounts/login/")
	req.Header.Set("Origin", "https://www.instagram.com")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Instagram-AJAX", c.rolloutHash)

	csrf := c.findCSRF(uri)

	if csrf != "" {
		req.Header.Set("X-CSRFToken", csrf)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return false, err
	}
	if resp.StatusCode != 200 {
		return false, fmt.Errorf("%d, %s", resp.StatusCode, resp.Status)
	} else {
		str, _ := ioutil.ReadAll(resp.Body)
		temp := map[string]interface{}{}
		json.Unmarshal([]byte(str), &temp)

		ok := temp["status"].(string) == "ok" && temp["authenticated"].(bool)

		if c.acc.Id == "" && ok {
			c.acc.Id = temp["userId"].(string)
		}
		return ok, nil
	}

	return false, err
}

func (c *Client) QueryComments(code string, cb func(data.Comment) (bool, error)) error {
	var err error
	var cursor = ""
	var res *structures.ShortCodeMediaResponse
	suspender := defaultSuspender{}
	for {
		res, err, cursor = c.GetShortCodeMediaInfo(code, cursor)
		if res != nil && res.Data.ShortCodeMedia != nil {
			for _, e := range res.Data.ShortCodeMedia.EdgeMediaToComment.Edges {
				comment := data.Comment{
					Id:        e.Node.ID,
					Text:      e.Node.Text,
					CreatedAt: e.Node.CreatedAt,
					Owner: owner.Owner{
						Id:            e.Node.Owner.ID,
						Username:      e.Node.Owner.Username,
						ProfilePicUrl: e.Node.Owner.ProfilePicURL,
					},
				}
				if r, err := cb(comment); !r {
					return err
				}
			}
		}
		if err != nil {
			if _, is := err.(httpErrors.EndOfListError); is {
				return nil
			}
			return err
		}
		suspender.Sleep()
	}
}

func (c *Client) QueryTag(tag string, cb func(data.TagMedia, *int) (bool, error)) error {
	var err error
	var cursor = ""
	var res *structures.HashTagResponse
	suspender := defaultSuspender{}
	for {
		res, err, cursor = c.GetTagPosts(tag, cursor)
		if res != nil && res.Data.HashTag != nil {
			var counter int = 0
			for _, e := range res.Data.HashTag.EdgeHashTagToMedia.Edges {
				post := data.TagMedia{
					Id:           e.Node.ID,
					Type:         e.Node.Typename,
					ShortCode:    e.Node.ShortCode,
					LikeCount:    int32(e.Node.EdgeLikedBy.Count),
					CommentCount: int32(e.Node.EdgeMediaToComment.Count),
					TakenAt:      e.Node.TakenAtTimestamp,
					Owner: owner.Owner{
						Id:       e.Node.Owner.ID,
						Username: "",
					},
				}
				if r, err := cb(post, &counter); !r {
					return err
				}
			}
			if counter == 0 {
				return nil
			}
		}
		if err != nil {
			if _, is := err.(httpErrors.EndOfListError); is {
				return nil
			}
			return err
		}
		suspender.Sleep()
	}
}

func (c *Client) SetAccount(acc *account.Account) {
	c.acc = acc
}

func (c *Client) GetAccount() *account.Account {
	return c.acc
}

func (c *Client) GetUserInfo(username string) (*data.User, error) {
	req, err := c.makeRequest("GET", fmt.Sprintf("https://www.instagram.com/%s/?__a=1", username), nil)
	if err != nil {
		return nil, err
	}
	c.prepare(req)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	str, _ := ioutil.ReadAll(resp.Body)
	temp := &structures.UserInfoResponse{}
	err = json.Unmarshal([]byte(str), &temp)
	if err != nil {
		return nil, fmt.Errorf("%v", temp)
	}
	u := &data.User{}
	u.Username = temp.GraphQL.User.Username
	u.Id = temp.GraphQL.User.ID
	u.Follows = int64(temp.GraphQL.User.EdgeFollow.Count)
	u.Followers = int64(temp.GraphQL.User.EdgeFollowedBy.Count)
	u.IsBusiness = temp.GraphQL.User.IsBusinessAccount
	u.IsPrivate = temp.GraphQL.User.IsPrivate
	u.IsVerified = temp.GraphQL.User.IsVerified
	return u, nil
}

func (c *Client) GetShortCodeMediaInfo(shortcode string, cursor string) (*structures.ShortCodeMediaResponse, error, string) {
	var query []byte
	var err error

	var hash string

	if cursor == "" {
		hash = "477b65a610463740ccdb83135b2014db"
		query, err = json.Marshal(map[string]interface{}{
			"shortcode":             shortcode,
			"child_comment_count":   0,
			"fetch_comment_count":   50,
			"parent_comment_count":  0,
			"has_threaded_comments": false,
		})
		if err != nil {
			return nil, err, ""
		}
	} else {
		hash = "f0986789a5c5d17c2400faebf16efd0d"
		query, err = json.Marshal(map[string]interface{}{
			"shortcode": shortcode,
			"first":     50,
			"after":     cursor,
		})
		if err != nil {
			return nil, err, ""
		}
	}

	req, err := c.makeRequest("GET", "https://www.instagram.com/graphql/query/", nil)
	q := req.URL.Query()
	q.Set("query_hash", hash)
	q.Set("variables", string(query))
	req.URL.RawQuery = q.Encode()
	c.prepare(req)

	resp, err := c.getHttpClient().Do(req)
	if err != nil {
		return nil, err, ""
	}
	if resp.StatusCode != 200 {
		switch resp.StatusCode {
		case 403:
			return nil, httpErrors.HttpForbidden{}, ""
		case 429:
			return nil, httpErrors.HttpTooManyRequests{}, ""
		default:
			return nil, errors.New(fmt.Sprintf("%d %s", resp.StatusCode, resp.Status)), ""
		}
	}
	bts, _ := ioutil.ReadAll(resp.Body)
	d := &structures.ShortCodeMediaResponse{}
	err = json.Unmarshal(bts, d)
	if err != nil {
		return nil, err, ""
	}
	if d.Data.ShortCodeMedia == nil {
		return nil, fmt.Errorf("%s", bts), ""
	}

	hasNext := d.Data.ShortCodeMedia.EdgeMediaToComment.PageInfo.HasNextPage
	if hasNext {
		cursor = d.Data.ShortCodeMedia.EdgeMediaToComment.PageInfo.EndCursor
	} else {
		return nil, httpErrors.EndOfListError{}, ""
	}

	return d, nil, cursor
}

func (c *Client) GetTagPosts(tag string, cursor string) (*structures.HashTagResponse, error, string) {
	var query []byte
	var err error

	hash := "f92f56d47dc7a55b606908374b43a314"

	variables := map[string]interface{}{
		"tag_name":    tag,
		"show_ranked": false,
		"first":       50,
	}

	if cursor != "" {
		variables["after"] = cursor
	}

	query, err = json.Marshal(variables)
	if err != nil {
		return nil, err, ""
	}

	req, _ := c.makeRequest("GET", "https://www.instagram.com/graphql/query/", nil)

	q := req.URL.Query()
	q.Set("query_hash", hash)

	q.Set("variables", string(query))
	req.URL.RawQuery = q.Encode()
	c.prepare(req)

	resp, err := c.getHttpClient().Do(req)
	if err != nil {
		return nil, err, ""
	}

	if resp.StatusCode != 200 {
		switch resp.StatusCode {
		case 403:
			return nil, httpErrors.HttpForbidden{}, ""
		case 429:
			return nil, httpErrors.HttpTooManyRequests{}, ""
		default:
			return nil, errors.New(fmt.Sprintf("%d %s", resp.StatusCode, resp.Status)), ""
		}
	}

	bts, _ := ioutil.ReadAll(resp.Body)
	d := &structures.HashTagResponse{}
	err = json.Unmarshal(bts, d)
	if err != nil {
		return nil, err, ""
	}

	hasNext := d.Data.HashTag.EdgeHashTagToMedia.PageInfo.HasNextPage
	if hasNext {
		cursor = d.Data.HashTag.EdgeHashTagToMedia.PageInfo.EndCursor
	} else {
		return nil, httpErrors.EndOfListError{}, ""
	}
	return d, nil, cursor
}

func (c *Client) GetShortCodeMediaLikers(shortcode string, cursor string) (*structures.ShortCodeMediaLikersResponse, error, string) {
	var query []byte
	var err error

	hash := "e0f59e4a1c8d78d0161873bc2ee7ec44"

	variables := map[string]interface{}{
		"shortcode":    shortcode,
		"include_reel": true,
		"first":        24,
	}

	if cursor != "" {
		variables["after"] = cursor
	}

	query, err = json.Marshal(variables)
	if err != nil {
		return nil, err, ""
	}

	req, _ := c.makeRequest("GET", "https://www.instagram.com/graphql/query/", nil)

	q := req.URL.Query()
	q.Set("query_hash", hash)

	q.Set("variables", string(query))
	req.URL.RawQuery = q.Encode()
	c.prepare(req)

	resp, err := c.getHttpClient().Do(req)
	if err != nil {
		return nil, err, ""
	}

	if resp.StatusCode != 200 {
		switch resp.StatusCode {
		case 403:
			return nil, httpErrors.HttpForbidden{}, ""
		case 429:
			return nil, httpErrors.HttpTooManyRequests{}, ""
		default:
			return nil, errors.New(fmt.Sprintf("%d %s", resp.StatusCode, resp.Status)), ""
		}
	}

	bts, _ := ioutil.ReadAll(resp.Body)
	d := &structures.ShortCodeMediaLikersResponse{}
	err = json.Unmarshal(bts, d)
	if err != nil {
		return nil, err, ""
	}

	hasNext := d.Data.ShortCodeMedia.EdgeLikedBy.PageInfo.HasNextPage
	if hasNext {
		cursor = d.Data.ShortCodeMedia.EdgeLikedBy.PageInfo.EndCursor
	} else {
		return nil, httpErrors.EndOfListError{}, ""
	}
	return d, nil, cursor
}

func (c *Client) GetHashTagSummary(tag string) (*structures.HashTagSummary, error) {
	req, err := c.makeRequest("GET", fmt.Sprintf("https://www.instagram.com/explore/tags/%s/", tag), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	respString, err := getResponseString(resp)
	if err != nil {
		return nil, err
	}

	sigStr := respString[strings.Index(respString, "_sharedData =")+13:]
	sigStr = sigStr[:strings.Index(sigStr, ";</script")]

	var res = &struct {
		EntryData struct {
			TagPage []struct {
				GraphQL struct {
					HashTag *structures.HashTagSummary `json:"hashtag"`
				} `json:"graphql"`
			} `json:"TagPage"`
		} `json:"entry_data"`
	}{}
	err = json.Unmarshal([]byte(sigStr), res)
	if err != nil {
		return nil, err
	}
	return res.EntryData.TagPage[0].GraphQL.HashTag, nil
}

func (c *Client) GetPostSummary(shortcode string) (*structures.PostSummary, error) {
	req, err := c.makeRequest("GET", fmt.Sprintf("https://www.instagram.com/p/%s/", shortcode), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	respString, err := getResponseString(resp)
	if err != nil {
		return nil, err
	}

	sigStr := respString[strings.Index(respString, "_sharedData =")+13:]
	sigStr = sigStr[:strings.Index(sigStr, ";</script")]

	var res = &struct {
		EntryData struct {
			PostPage []struct {
				GraphQL struct {
					ShortCodeMedia *structures.PostSummary `json:"shortcode_media"`
				} `json:"graphql"`
			} `json:"PostPage"`
		} `json:"entry_data"`
	}{}
	err = json.Unmarshal([]byte(sigStr), res)
	if err != nil {
		return nil, err
	}
	return res.EntryData.PostPage[0].GraphQL.ShortCodeMedia, nil
}
