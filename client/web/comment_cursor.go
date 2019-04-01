package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"giveaway/client"
	"giveaway/data"
	httpErrors "giveaway/data/errors"
	"io/ioutil"
)

type CommentCursor struct {
	hashes    [2]string
	step      int
	cursor    string
	hasNext   bool
	code      string
	client    *Client
	suspender client.SuspendsThread
}

func (c *CommentCursor) Next() (<-chan data.Comment, <-chan error) {
	resChan := make(chan data.Comment)
	errChan := make(chan error)

	go c.run(errChan, resChan)

	return resChan, errChan
}

func (c *CommentCursor) run(errChan chan error, resChan chan data.Comment) {
	var query []byte
	var err error
	var hash string

	finalizer := func() {
		close(errChan)
		close(resChan)
	}

	for {
		hash = c.hashes[c.step]

		if c.step == 0 {
			query, err = json.Marshal(map[string]interface{}{
				"shortcode":             c.code,
				"child_comment_count":   0,
				"fetch_comment_count":   50,
				"parent_comment_count":  0,
				"has_threaded_comments": false,
			})
			if err != nil {
				errChan <- err
				finalizer()
				return
			}
			c.step++
		} else {
			if !c.hasNext {
				finalizer()
				return
			}
			query, err = json.Marshal(map[string]interface{}{
				"shortcode": c.code,
				"first":     50,
				"after":     c.cursor,
			})
			if err != nil {
				errChan <- err
				finalizer()
				return
			}
		}

		req, err := c.client.makeRequest("GET", "https://www.instagram.com/graphql/query/", nil)
		q := req.URL.Query()
		q.Set("query_hash", hash)
		q.Set("variables", string(query))
		req.URL.RawQuery = q.Encode()
		c.client.prepare(req)

		resp, err := c.client.getHttpClient().Do(req)
		if err != nil {
			errChan <- err
			finalizer()
			return
		}
		if resp.StatusCode != 200 {
			switch resp.StatusCode {
			case 403:
				errChan <- httpErrors.HttpForbidden{}
			case 429:
				errChan <- httpErrors.HttpTooManyRequests{}
			default:
				errChan <- errors.New(fmt.Sprintf("%d %s", resp.StatusCode, resp.Status))
			}
			finalizer()
			return
		}
		bytes, _ := ioutil.ReadAll(resp.Body)
		d := map[string]interface{}{}
		err = json.Unmarshal(bytes, &d)
		if err != nil {
			errChan <- err
			finalizer()
			return
		}
		temporaryObject := d["data"].(map[string]interface{})["shortcode_media"].(map[string]interface{})["edge_media_to_comment"].(map[string]interface{})
		pageInfo := temporaryObject["page_info"].(map[string]interface{})
		c.hasNext = pageInfo["has_next_page"].(bool)
		if c.hasNext {
			c.cursor = pageInfo["end_cursor"].(string)
		}
		edges := temporaryObject["edges"].([]interface{})

		for _, node := range edges {
			temp := node.(map[string]interface{})["node"].(map[string]interface{})
			if owner, p := temp["owner"].(map[string]interface{}); p {
				resChan <- data.Comment{
					Id:        temp["id"].(string),
					Text:      temp["text"].(string),
					CreatedAt: int64(temp["created_at"].(float64)),
					Owner: data.Owner{
						Id:       owner["id"].(string),
						Username: owner["username"].(string),
					},
				}
			}
		}
		c.suspender.Sleep()
	}
}

func (c *CommentCursor) SetSuspender(suspender client.SuspendsThread) {
	c.suspender = suspender
}

func NewCommentCursor(c *Client, code string) *CommentCursor {
	cursor := &CommentCursor{}
	cursor.step = 0
	cursor.cursor = ""
	cursor.hasNext = false
	cursor.code = code
	cursor.client = c
	cursor.hashes = [2]string{
		"477b65a610463740ccdb83135b2014db",
		"f0986789a5c5d17c2400faebf16efd0d",
	}

	cursor.suspender = &defaultSuspender{}

	return cursor
}
