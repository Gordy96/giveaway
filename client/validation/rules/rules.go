package rules

import (
	"encoding/json"
	"fmt"
	"giveaway/client"
	"giveaway/client/api"
	"giveaway/client/web"
	"giveaway/data/errors"
	"giveaway/data/сontainers"
	"giveaway/instagram"
	"giveaway/instagram/account"
	"giveaway/instagram/account/repository"
	"giveaway/instagram/commands"
	"giveaway/instagram/solver"
	"giveaway/instagram/structures/stories"
	"giveaway/utils/bson"
	"time"
)

type DateRule struct {
	Name   string   `json:"name" bson:"name"`
	Limits [2]int64 `json:"limits" bson:"limits"`
}

func (d DateRule) String() string {
	bts, _ := json.Marshal(d)
	return string(bts)
}

func (d DateRule) MarshalBSON() ([]byte, error) {
	return bson.StructToBSON(d)
}

func (d DateRule) UnmarshalBSON(data []byte) error {
	return bson.BSONToStruct(data, &d)
}

func (d DateRule) Validate(i interface{}) (bool, error) {
	e, is := i.(client.HasDateAttribute)
	if !is {
		return false, errors.ValidationCriticalFailure{}
	}
	examined := e.GetCreationDate()
	if d.Limits[1] > 0 && examined > d.Limits[1] {
		return false, errors.AfterMaximumDate{}
	}
	if examined < d.Limits[0] {
		return false, errors.BeforeMinimumDate{}
	}
	return true, nil
}

type FollowsRule struct {
	Name string `json:"name" bson:"name"`
	Id   string `json:"id" bson:"id"`
}

func (f FollowsRule) String() string {
	bts, _ := json.Marshal(f)
	return string(bts)
}

func (f FollowsRule) MarshalBSON() ([]byte, error) {
	return bson.StructToBSON(f)
}

func (f FollowsRule) UnmarshalBSON(data []byte) error {
	return bson.BSONToStruct(data, &f)
}

func (f FollowsRule) Validate(i interface{}) (bool, error) {
	var is bool
	var err error = nil
	var hasOwner client.HasOwner
	hasOwner, is = i.(client.HasOwner)
	if !is {
		return false, errors.ValidationCriticalFailure{}
	}
	owner := hasOwner.GetOwner()
	repo := repository.GetRepositoryInstance()
	cl := api.NewApiClient()

	for {
		acc := repo.GetOldestUsedRetries(15, 2*time.Second)
		if acc == nil {
			return false, errors.ValidationCriticalFailure{}
		}
		cl.SetAccount(acc)
		is, err = cl.IsFollower(owner, f.Id)
		if err == nil {
			acc.Status = account.Available
			repo.Save(acc)
			break
		}
		switch err.(type) {
		case errors.LoginRequired:
			solver.GetRunningInstance().Enqueue(commands.MakeNewReLoginCommand(acc))
		}
	}
	return is, err
}

func checkCondition(con string, arg0, arg1 int64) bool {
	switch con {
	case ">":
		return arg0 > arg1
	case ">=":
		return arg0 >= arg1
	case "<":
		return arg0 < arg1
	case "<=":
		return arg0 <= arg1
	case "==":
		return arg0 == arg1
	case "!=":
		return arg0 != arg1
	}
	return false
}

type FollowersRule struct {
	Name      string `json:"name" bson:"name"`
	Amount    int64  `json:"amount" bson:"amount"`
	Username  string `json:"username" bson:"username"`
	Condition string `json:"condition" bson:"condition"`
}

func (f FollowersRule) MarshalBSON() ([]byte, error) {
	return bson.StructToBSON(f)
}

func (f FollowersRule) UnmarshalBSON(data []byte) error {
	return bson.BSONToStruct(data, &f)
}

func (f FollowersRule) String() string {
	bts, _ := json.Marshal(f)
	return string(bts)
}

func (f FollowersRule) Validate(i interface{}) (bool, error) {
	cl, is := i.(*web.Client)
	if !is {
		return false, errors.ValidationCriticalFailure{}
	}
	u, err := cl.GetUserInfo(f.Username)
	if err != nil {
		return false, err
	}
	return checkCondition(f.Condition, u.Followers, f.Amount), nil
}

type PostLikesRule struct {
	Name      string `json:"name" bson:"name"`
	Amount    int64  `json:"amount" bson:"amount"`
	ShortCode string `json:"shortcode" bson:"shortcode"`
	Condition string `json:"condition" bson:"condition"`
}

func (p PostLikesRule) MarshalBSON() ([]byte, error) {
	return bson.StructToBSON(p)
}

func (p PostLikesRule) UnmarshalBSON(data []byte) error {
	return bson.BSONToStruct(data, &p)
}

func (p PostLikesRule) String() string {
	bts, _ := json.Marshal(p)
	return string(bts)
}

func (p PostLikesRule) Validate(i interface{}) (bool, error) {
	cl := i.(*web.Client)
	u, err, _ := cl.GetShortCodeMediaLikers(p.ShortCode, "")
	if err != nil {
		return false, err
	}
	return checkCondition(p.Condition, u.Data.ShortCodeMedia.EdgeLikedBy.Count, p.Amount), nil
}

type PostCommentsRule struct {
	Name      string `json:"name" bson:"name"`
	Amount    int64  `json:"amount" bson:"amount"`
	ShortCode string `json:"shortcode" bson:"shortcode"`
	Condition string `json:"condition" bson:"condition"`
}

func (p PostCommentsRule) MarshalBSON() ([]byte, error) {
	return bson.StructToBSON(p)
}

func (p PostCommentsRule) UnmarshalBSON(data []byte) error {
	return bson.BSONToStruct(data, &p)
}

func (p PostCommentsRule) String() string {
	bts, _ := json.Marshal(p)
	return string(bts)
}

func (p PostCommentsRule) Validate(i interface{}) (bool, error) {
	cl := i.(*web.Client)
	u, err, _ := cl.GetShortCodeMediaInfo(p.ShortCode, "")
	if err != nil {
		return false, err
	}
	return checkCondition(p.Condition, u.Data.ShortCodeMedia.EdgeMediaToComment.Count, p.Amount), nil
}

type ParticipantsRule struct {
	Name      string `json:"name" bson:"name"`
	Amount    int64  `json:"amount" bson:"amount"`
	Condition string `json:"condition" bson:"condition"`
}

func (p ParticipantsRule) MarshalBSON() ([]byte, error) {
	return bson.StructToBSON(p)
}

func (p ParticipantsRule) UnmarshalBSON(data []byte) error {
	return bson.BSONToStruct(data, &p)
}

func (p ParticipantsRule) String() string {
	bts, _ := json.Marshal(p)
	return string(bts)
}

func (p ParticipantsRule) Validate(i interface{}) (bool, error) {
	col, ok := i.(*сontainers.EntryContainer)
	if !ok {
		return false, fmt.Errorf("wrong argumane type (%v)", i)
	}
	return checkCondition(p.Condition, int64(col.LengthNoDuplicates()), p.Amount), nil
}

type StoryHasHashTagRule struct {
	Name    string `json:"name" bson:"name"`
	HashTag string `json:"hashtag" bson:"hashtag"`
}

func (s StoryHasHashTagRule) MarshalBSON() ([]byte, error) {
	return bson.StructToBSON(s)
}

func (s StoryHasHashTagRule) UnmarshalBSON(data []byte) error {
	return bson.BSONToStruct(data, &s)
}

func (s StoryHasHashTagRule) String() string {
	bts, _ := json.Marshal(s)
	return string(bts)
}

func (s StoryHasHashTagRule) Validate(i interface{}) (bool, error) {
	story, ok := i.(*stories.StoryItem)
	if !ok {
		return false, fmt.Errorf("wrong argumane type (%v)", i)
	}

	for _, tag := range story.StoryHashtags {
		if tag.Hashtag.Name == s.HashTag {
			return true, nil
		}
	}
	return false, nil
}

type StoryHasMentionRule struct {
	Name     string `json:"name" bson:"name"`
	Username string `json:"username, omitempty" bson:"username, omitempty"`
	ID       string `json:"id, omitempty" bson:"id, omitempty"`
}

func (s StoryHasMentionRule) MarshalBSON() ([]byte, error) {
	return bson.StructToBSON(s)
}

func (s StoryHasMentionRule) UnmarshalBSON(data []byte) error {
	return bson.BSONToStruct(data, &s)
}

func (s StoryHasMentionRule) String() string {
	bts, _ := json.Marshal(s)
	return string(bts)
}

func (s StoryHasMentionRule) Validate(i interface{}) (bool, error) {
	story, ok := i.(*stories.StoryItem)
	if !ok {
		return false, fmt.Errorf("wrong argumane type (%v)", i)
	}

	for _, mention := range story.ReelMentions {
		if s.Username != "" {
			if mention.User.Username == s.Username {
				return true, nil
			}
		} else if s.ID != "" {
			if string(mention.User.Pk) == s.ID {
				return true, nil
			}
		} else {
			return false, nil
		}
	}
	return false, nil
}

type StoryHasPostRule struct {
	Name      string `json:"name" bson:"name"`
	ShortCode string `json:"shortcode, omitempty" bson:"shortcode, omitempty"`
}

func (s StoryHasPostRule) MarshalBSON() ([]byte, error) {
	return bson.StructToBSON(s)
}

func (s StoryHasPostRule) UnmarshalBSON(data []byte) error {
	return bson.BSONToStruct(data, &s)
}

func (s StoryHasPostRule) String() string {
	bts, _ := json.Marshal(s)
	return string(bts)
}

func (s StoryHasPostRule) Validate(i interface{}) (bool, error) {
	story, ok := i.(*stories.StoryItem)
	if !ok {
		return false, fmt.Errorf("wrong argumane type (%v)", i)
	}

	for _, media := range story.StoryFeedMedia {
		if instagram.IdToCode(media.MediaID) == s.ShortCode {
			return true, nil
		}
	}
	return false, nil
}

type StoryHasExternalLinkRule struct {
	Name string `json:"name" bson:"name"`
	Link string `json:"link, omitempty" bson:"link, omitempty"`
}

func (s StoryHasExternalLinkRule) MarshalBSON() ([]byte, error) {
	return bson.StructToBSON(s)
}

func (s StoryHasExternalLinkRule) UnmarshalBSON(data []byte) error {
	return bson.BSONToStruct(data, &s)
}

func (s StoryHasExternalLinkRule) String() string {
	bts, _ := json.Marshal(s)
	return string(bts)
}

func (s StoryHasExternalLinkRule) Validate(i interface{}) (bool, error) {
	story, ok := i.(*stories.StoryItem)
	if !ok {
		return false, fmt.Errorf("wrong argumane type (%v)", i)
	}

	for _, cta := range story.StoryCta {
		for _, link := range cta.Links {
			if link.WebURI == s.Link {
				return true, nil
			}
		}
	}
	return false, nil
}
