package requests

import "giveaway/client/validation"

type HasRulesJsonRequest struct {
	Rules validation.RuleCollection `json:"rules"`
}

type CommentTaskJsonRequest struct {
	HasRulesJsonRequest `json:",inline" bson:",inline"`
	ShortCode           string `json:"shortcode" bson:"shortcode"`
}

type HashTagTaskJsonRequest struct {
	HasRulesJsonRequest `json:",inline" bson:",inline"`
	HashTag             string `json:"hashtag" bson:"hashtag"`
}