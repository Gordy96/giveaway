package data

import (
	"giveaway/client/validation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Owner struct {
	Id       string `json:"id" bson:"id"`
	Username string `json:"username" bson:"username"`
}

type User struct {
	Owner      `json:",inline" bson:",inline"`
	Follows    int64 `json:"follows" bson:"follows"`
	Followers  int64 `json:"followers" bson:"followers"`
	IsBusiness bool  `json:"is_business" bson:"is_business"`
	IsPrivate  bool  `json:"is_private" bson:"is_private"`
	IsVerified bool  `json:"is_verified" bson:"is_verified"`
}

type Comment struct {
	Id        string `json:"id,omitempty" bson:"id"`
	Text      string `json:"text" bson:"text"`
	Owner     Owner  `json:"owner" bson:"owner"`
	CreatedAt int64  `json:"created_at" bson:"created_at"`
}

func (c *Comment) GetOwner() *Owner {
	return &c.Owner
}

func (c *Comment) GetKey() interface{} {
	return c.Owner.Id
}

func (c *Comment) GetCreationDate() int64 {
	return c.CreatedAt
}

type TagMedia struct {
	Id           string `json:"id" bson:"id"`
	Type         string `json:"type" bson:"type"`
	ShortCode    string `json:"shortcode" bson:"shortcode"`
	LikeCount    int32  `json:"like_count" bson:"like_count"`
	CommentCount int32  `json:"comment_count" bson:"comment_count"`
	TakenAt      int32  `json:"taken_at" bson:"taken_at"`
	Owner        Owner  `json:"owner" bson:"owner"`
}

func (t *TagMedia) GetCreationDate() int32 {
	return t.TakenAt
}

func (t *TagMedia) GetOwner() *Owner {
	return &t.Owner
}

func (t *TagMedia) GetKey() interface{} {
	return t.Owner.Id
}

type BaseTaskModel struct {
	Id                 primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	CreatedAt          int64              `json:"created_at" bson:"created_at"`
	FinishedAt         int64              `json:"finished_at" bson:"finished_at"`
	Status             string             `json:"status" bson:"status"`
	Comment            string             `json:"comment" bson:"comment"`
	PreconditionRules  []validation.IRule `json:"precondition_rules" bson:"precondition_rules"`
	AppendRules        []validation.IRule `json:"append_rules" append_rules:"rules"`
	SelectRules        []validation.IRule `json:"select_rules" bson:"select_rules"`
	PostconditionRules []validation.IRule `json:"postcondition_rules" bson:"postcondition_rules"`
}

type HasKey interface {
	GetKey() interface{}
}

type CommentsTask struct {
	BaseTaskModel `json:",inline" bson:",inline"`
	ShortCode     string     `json:"shortcode" bson:"shortcode"`
	Winner        *Comment   `json:"winner" bson:"winner"`
	Above         []*Comment `json:"above" bson:"above"`
	Below         []*Comment `json:"below" bson:"below"`
	Position      int        `json:"position"`
}

func (c *CommentsTask) GetKey() interface{} {
	return c.Id
}

type HashTagTask struct {
	BaseTaskModel `json:",inline" bson:",inline"`
	HashTag       string    `json:"hashtag" bson:"hashtag"`
	Post          *TagMedia `json:"post,omitempty" bson:"post"`
}

func (c *HashTagTask) GetKey() interface{} {
	return c.Id
}
