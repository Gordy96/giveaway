package web

import (
	"giveaway/client"
	"giveaway/data"
	"giveaway/data/owner"
	"giveaway/instagram/structures"
)

type CommentCursor struct {
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
	defer func() {
		close(errChan)
		close(resChan)
	}()
	var err error
	var cursor = ""
	var res *structures.ShortCodeMediaResponse
	for {
		res, err, cursor = c.client.GetShortCodeMediaInfo(c.code, cursor)
		if err != nil {
			errChan <- err
			return
		}

		for _, e := range res.Data.ShortCodeMedia.EdgeMediaToComment.Edges {
			resChan <- data.Comment{
				Id:        e.Node.ID,
				Text:      e.Node.Text,
				CreatedAt: e.Node.CreatedAt,
				Owner: owner.Owner{
					Id:       e.Node.Owner.ID,
					Username: e.Node.Owner.Username,
				},
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
	cursor.client = c
	cursor.code = code
	cursor.suspender = &defaultSuspender{}

	return cursor
}
