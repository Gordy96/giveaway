package data

import "go.mongodb.org/mongo-driver/bson/primitive"

type Owner struct {
	Id 				string					`json:"id" bson:"id"`
	Username 		string					`json:"username" bson:"username"`
}

type Comment struct {
	Id 				string					`json:"id,omitempty" bson:"id"`
	Text 			string					`json:"text" bson:"text"`
	Owner 			Owner					`json:"owner" bson:"owner"`
	CreatedAt		int32					`json:"created_at" bson:"created_at"`
}

func (t *Comment) GetCreationDate() int32 {
	return t.CreatedAt
}

type TagMedia struct {
	Id 				string					`json:"id" bson:"id"`
	Type 			string					`json:"type" bson:"type"`
	ShortCode 		string					`json:"shortcode" bson:"shortcode"`
	LikeCount		int32					`json:"like_count" bson:"like_count"`
	CommentCount	int32					`json:"comment_count" bson:"comment_count"`
	TakenAt 		int32					`json:"taken_at" bson:"taken_at"`
	Owner 			Owner					`json:"owner" bson:"owner"`
}

func (t *TagMedia) GetCreationDate() int32 {
	return t.TakenAt
}

type BaseTaskModel struct {
	Id				primitive.ObjectID		`json:"_id,omitempty" bson:"_id"`
	CreatedAt		int32					`json:"created_at" bson:"created_at"`
	FinishedAt		int32					`json:"finished_at" bson:"finished_at"`
}

type CommentsTask struct {
	BaseTaskModel							`json:",inline" bson:",inline"`
	ShortCode		string					`json:"shortcode" bson:"shortcode"`
	Winner			*Comment				`json:"winner" bson:"winner"`
	Above			[]*Comment				`json:"above" bson:"above"`
	Below			[]*Comment				`json:"below" bson:"below"`
	Status 			string					`json:"status" bson:"status"`
	Position		int						`json:"position"`
}

type HashTagTask struct {
	BaseTaskModel							`json:",inline" bson:",inline"`
	HashTag			string					`json:"hashtag" bson:"hashtag"`
	Post			*TagMedia				`json:"post,omitempty" bson:"post"`
	Status 			string					`json:"status" bson:"status"`
}