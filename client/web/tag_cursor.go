package web

import (
	"giveaway/client"
	"giveaway/data"
	"giveaway/instagram/structures"
)

type TagMediaCursor struct {
	cursor    string
	hasNext   bool
	tag       string
	client    *Client
	suspender client.SuspendsThread
}

func (c *TagMediaCursor) Next() (<-chan data.TagMedia, <-chan error) {
	resChan := make(chan data.TagMedia)
	errChan := make(chan error)

	go c.run(errChan, resChan)

	return resChan, errChan
}

func (c *TagMediaCursor) run(errChan chan error, resChan chan data.TagMedia) {
	defer func() {
		close(errChan)
		close(resChan)
	}()
	var err error
	var cursor = ""
	var res *structures.HashTagResponse
	for {
		res, err, cursor = c.client.GetTagPosts(c.tag, cursor)
		if err != nil {
			errChan <- err
			return
		}

		for _, e := range res.Data.HashTag.EdgeHashTagToMedia.Edges {
			resChan <- data.TagMedia{
				Id:           e.Node.ID,
				Type:         e.Node.Typename,
				ShortCode:    e.Node.ShortCode,
				LikeCount:    int32(e.Node.EdgeLikedBy.Count),
				CommentCount: int32(e.Node.EdgeMediaToComment.Count),
				TakenAt:      int32(e.Node.TakenAtTimestamp),
				Owner: data.Owner{
					Id:       e.Node.Owner.ID,
					Username: "",
				},
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
