package data

import (
	"giveaway/data/owner"
	"giveaway/instagram/structures/stories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	owner.Owner `json:",inline" bson:",inline"`
	Follows     int64 `json:"follows" bson:"follows"`
	Followers   int64 `json:"followers" bson:"followers"`
	IsBusiness  bool  `json:"is_business" bson:"is_business"`
	IsPrivate   bool  `json:"is_private" bson:"is_private"`
	IsVerified  bool  `json:"is_verified" bson:"is_verified"`
}

type Comment struct {
	Id        string      `json:"id,omitempty" bson:"id"`
	Text      string      `json:"text" bson:"text"`
	Owner     owner.Owner `json:"owner" bson:"owner"`
	CreatedAt int64       `json:"created_at" bson:"created_at"`
}

func (c *Comment) GetOwner() *owner.Owner {
	return &c.Owner
}

func (c *Comment) GetKey() interface{} {
	return c.Owner.Id
}

func (c *Comment) GetCreationDate() int64 {
	return c.CreatedAt
}

type TagMedia struct {
	Id           string      `json:"id" bson:"id"`
	Type         string      `json:"type" bson:"type"`
	ShortCode    string      `json:"shortcode" bson:"shortcode"`
	LikeCount    int32       `json:"like_count" bson:"like_count"`
	CommentCount int32       `json:"comment_count" bson:"comment_count"`
	TakenAt      int64       `json:"taken_at" bson:"taken_at"`
	Owner        owner.Owner `json:"owner" bson:"owner"`
}

func (t *TagMedia) GetCreationDate() int64 {
	return t.TakenAt
}

func (t *TagMedia) GetOwner() *owner.Owner {
	return &t.Owner
}

func (t *TagMedia) GetKey() interface{} {
	return t.Owner.Id
}

type CommentWinner struct {
	Winner   *Comment   `json:"winner" bson:"winner"`
	Above    []*Comment `json:"above" bson:"above"`
	Below    []*Comment `json:"below" bson:"below"`
	Position int        `json:"position" bson:"position"`
}

type CommentContainer struct {
	Id      primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	TaskId  primitive.ObjectID `json:"task_id,omitempty" bson:"task_id"`
	Comment Comment            `json:"comment" bson:"comment"`
}

func (c CommentContainer) GetKey() interface{} {
	return c.Id
}

type PostContainer struct {
	Id     primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	TaskId primitive.ObjectID `json:"task_id,omitempty" bson:"task_id"`
	Post   TagMedia           `json:"post" bson:"post"`
}

func (p PostContainer) GetKey() interface{} {
	return p.Id
}

type StoryContainer struct {
	Id        int64              `json:"_id,omitempty" bson:"_id"`
	CreatedAt int64              `json:"created_at" bson:"created_at"`
	TaskId    primitive.ObjectID `json:"task_id,omitempty" bson:"task_id"`
	Story     stories.StoryItem  `json:"story" bson:"story"`
	Data      []byte             `json:"data" bson:"data"`
}

func (p StoryContainer) GetKey() interface{} {
	return p.Id
}
