package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"giveaway/client"
	"giveaway/data"
	httpErrors "giveaway/data/errors"
	"io/ioutil"
	"time"
)

type TagMediaCursor struct {
	cursor 		string
	hasNext 	bool
	tag 		string
	client 		*Client
	suspender	client.SuspendsThread
}

func (c *TagMediaCursor) Next() <-chan data.TagMedia {
	resChan := make(chan data.TagMedia)
	errChan := make(chan error)
	go func() {
		for {
			e, ok := <- errChan
			if !ok {
				return
			}
			fmt.Printf("[%s]: %s", time.Now().String(), e.Error())
		}
	}()

	go c.run(errChan, resChan)

	return resChan
}

func (c *TagMediaCursor) run(errChan chan error, resChan chan data.TagMedia){
	var query []byte
	var err error

	hash := "f92f56d47dc7a55b606908374b43a314"

	finalizer := func() {
		close(errChan)
		close(resChan)
	}
	c.hasNext = true
	for {
		if !c.hasNext {
			finalizer()
			return
		}

		variables := map[string]interface{} {
			"tag_name": c.tag,
			"show_ranked":false,
			"first": 50,
		}

		if c.cursor != "" {
			variables["after"] = c.cursor
		}

		query, err = json.Marshal(variables)
		if err != nil {
			errChan <- err
			finalizer()
			return
		}

		req, _ := c.client.makeRequest("GET", "https://www.instagram.com/graphql/query/", nil)

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
		}
		temporaryObject := d["data"].(map[string]interface{})["hashtag"].(map[string]interface{})["edge_hashtag_to_media"].(map[string]interface{})
		pageInfo := temporaryObject["page_info"].(map[string]interface{})
		c.hasNext = pageInfo["has_next_page"].(bool)
		if c.hasNext {
			c.cursor = pageInfo["end_cursor"].(string)
		}
		edges := temporaryObject["edges"].([]interface{})

		for _, node := range edges {
			temp := node.(map[string]interface{})["node"].(map[string]interface{})
			if owner, p := temp["owner"].(map[string]interface{}); p {
				liked, _ := temp["edge_liked_by"].(map[string]interface{})["count"].(float64)
				commented, _ := temp["edge_media_to_comment"].(map[string]interface{})["count"].(float64)
				timestamp, _ := temp["taken_at_timestamp"].(float64)
				resChan <- data.TagMedia{
					Id: temp["id"].(string),
					Type: temp["__typename"].(string),
					ShortCode: temp["shortcode"].(string),
					LikeCount: int32(liked),
					CommentCount: int32(commented),
					TakenAt: int32(timestamp),
					Owner: data.Owner{
						Id: owner["id"].(string),
						Username: "",
					},
				}
			}
		}
		c.suspender.Sleep()
	}

}

func (c *TagMediaCursor) SetSuspender(suspender client.SuspendsThread) {
	c.suspender = suspender
}

func NewTagCursor(c *Client, tag string) *TagMediaCursor {
	cursor := &TagMediaCursor{}
	cursor.cursor = ""
	cursor.hasNext = true
	cursor.tag = tag
	cursor.client = c
	return cursor
}