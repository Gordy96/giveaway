package web

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"giveaway/data"
	"giveaway/instagram/account"
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
	cursor := NewCommentCursor(c, code)
	resChan, errChan := cursor.Next()
	for {
		select {
		case post := <-resChan:
			if r, err := cb(post); !r {
				return err
			}
		case err := <-errChan:
			return err
		}
	}
}

func (c *Client) QueryTag(tag string, cb func(data.TagMedia) (bool, error)) error {
	cursor := NewTagCursor(c, tag)
	resChan, errChan := cursor.Next()
	for {
		select {
		case post := <-resChan:
			if r, err := cb(post); !r {
				return err
			}
		case err := <-errChan:
			return err
		}
	}
}

func (c *Client) SetAccount(acc *account.Account) {
	c.acc = acc
}

func (c *Client) GetAccount() *account.Account {
	return c.acc
}
