package data

import "go.mongodb.org/mongo-driver/bson/primitive"

type Owner struct {
	Id       string `json:"id" bson:"id"`
	Username string `json:"username" bson:"username"`
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

type BaseTaskModel struct {
	Id         primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	CreatedAt  int64              `json:"created_at" bson:"created_at"`
	FinishedAt int64              `json:"finished_at" bson:"finished_at"`
}

type DatabaseTaskResult interface {
	CommentsTask() *CommentsTask
	HashTagTask() *HashTagTask
}

func NewSingleDocumentResult(i interface{}) *SingleDocumentResult {
	return &SingleDocumentResult{i}
}

type SingleDocumentResult struct {
	value interface{}
}

func (s *SingleDocumentResult) CommentsTask() *CommentsTask {
	return s.value.(*CommentsTask)
}

func (s *SingleDocumentResult) HashTagTask() *HashTagTask {
	return s.value.(*HashTagTask)
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
	Status        string     `json:"status" bson:"status"`
	Position      int        `json:"position"`
}

func (c *CommentsTask) GetKey() interface{} {
	return c.Id
}

type HashTagTask struct {
	BaseTaskModel `json:",inline" bson:",inline"`
	HashTag       string    `json:"hashtag" bson:"hashtag"`
	Post          *TagMedia `json:"post,omitempty" bson:"post"`
	Status        string    `json:"status" bson:"status"`
}

func (c *HashTagTask) GetKey() interface{} {
	return c.Id
}
