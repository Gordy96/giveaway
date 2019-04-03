package rules

import (
	"encoding/json"
	"fmt"
	"giveaway/client"
	"giveaway/client/api"
	"giveaway/client/web"
	"giveaway/data/errors"
	"giveaway/data/сontainers"
	"giveaway/instagram/account"
	"giveaway/instagram/account/repository"
	"giveaway/instagram/commands"
	"giveaway/instagram/solver"
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
	examined := i.(client.HasDateAttribute).GetCreationDate()
	if d.Limits[1] > 0 && examined > d.Limits[1] {
		return false, nil
	}
	if examined < d.Limits[0] {
		return false, errors.ShouldStopIterationError{}
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
	owner := i.(client.HasOwner).GetOwner()
	repo := repository.GetRepositoryInstance()
	cl := api.NewApiClient()

	var is bool
	var err error = nil

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
	cl := i.(*web.Client)
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
