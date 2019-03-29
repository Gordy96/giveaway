package rules

import (
	"giveaway/client"
	"giveaway/client/api"
	"giveaway/data/errors"
	"giveaway/instagram/account/repository"
	"giveaway/instagram/commands"
	"giveaway/instagram/solver"
)

type DateRule struct {
	Name   string   `json:"name"`
	Limits [2]int64 `json:"limits"`
}

func (d DateRule) GetName() string {
	return d.Name
}

func (d DateRule) Validate(i interface{}) (bool, error) {
	examined := i.(client.HasDateAttribute).GetCreationDate()
	if examined > d.Limits[1] {
		return false, nil
	}
	if examined < d.Limits[0] {
		return false, errors.ShouldStopIterationError{}
	}
	return true, nil
}

type FollowingRule struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

func (f FollowingRule) GetName() string {
	return f.Name
}

func (f FollowingRule) Validate(i interface{}) (bool, error) {
	owner := i.(client.HasOwner).GetOwner()
	repo := repository.GetRepositoryInstance()
	cl := api.NewApiClient()

	var is bool
	var err error = nil

	for {
		acc := repo.GetOldestUsed()
		if acc == nil {
			return false, errors.ValidationCriticalFailure{}
		}
		cl.SetAccount(acc)
		is, err = cl.IsFollower(owner, f.Id)
		if err == nil {
			break
		}
		switch err.(type) {
		case errors.LoginRequired:
			solver.GetRunningInstance().Enqueue(commands.MakeNewReLoginCommand(acc))
		}

	}
	return is, err
}
