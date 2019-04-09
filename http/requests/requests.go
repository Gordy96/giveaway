package requests

import (
	"giveaway/client/validation"
)

type HasRulesJsonRequest struct {
	Rules validation.RuleCollection `json:"rules"`
}

type HasWinnerCount struct {
	NumWinners int `json:"num_winners"`
}

type CommentTaskJsonRequest struct {
	HasRulesJsonRequest `json:",inline" bson:",inline"`
	HasWinnerCount      `json:",inline" bson:",inline"`
	ShortCode           string `json:"shortcode" bson:"shortcode"`
}

type HashTagTaskJsonRequest struct {
	HasRulesJsonRequest `json:",inline" bson:",inline"`
	HasWinnerCount      `json:",inline" bson:",inline"`
	HashTag             string `json:"hashtag" bson:"hashtag"`
}

type HashTagStoryTaskJsonRequest struct {
	HasRulesJsonRequest `json:",inline" bson:",inline"`
	HasWinnerCount      `json:",inline" bson:",inline"`
	HashTag             string `json:"hashtag" bson:"hashtag"`
}
