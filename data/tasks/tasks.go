package tasks

import (
	"giveaway/client/validation"
	"giveaway/data"
	"giveaway/instagram/structures"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseTaskModel struct {
	Id         primitive.ObjectID        `json:"_id,omitempty" bson:"_id"`
	CreatedAt  int64                     `json:"created_at" bson:"created_at"`
	FinishedAt int64                     `json:"finished_at" bson:"finished_at"`
	Status     string                    `json:"status" bson:"status"`
	Comment    string                    `json:"comment" bson:"comment"`
	Rules      validation.RuleCollection `json:"rules" bson:"rules"`
	NumWinners int                       `json:"num_winners"`
}

type CommentsTask struct {
	BaseTaskModel `json:",inline" bson:",inline"`
	ShortCode     string               `json:"shortcode" bson:"shortcode"`
	Winners       []data.CommentWinner `json:"winners" bson:"winners"`
}

func (c *CommentsTask) GetKey() interface{} {
	return c.Id
}

type HashTagTask struct {
	BaseTaskModel `json:",inline" bson:",inline"`
	HashTag       string           `json:"hashtag" bson:"hashtag"`
	Winners       []*data.TagMedia `json:"winners" bson:"winners"`
}

func (c *HashTagTask) GetKey() interface{} {
	return c.Id
}

type StoriesTask struct {
	BaseTaskModel `json:",inline" bson:",inline"`
	HashTag       string                 `json:"hashtag" bson:"hashtag"`
	Winners       []structures.StoryItem `json:"winners" bson:"winners"`
	Account       string                 `json:"account" bson:"account"`
}

func (c *StoriesTask) GetKey() interface{} {
	return c.Id
}
