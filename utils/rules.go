package utils

import (
	"encoding/json"
	"giveaway/client"
	"giveaway/data/errors"
)

type RuleType int

const (
	AppendRule RuleType = iota
	SelectRule RuleType = AppendRule + 1
)

type ConstructorFunc func(interface{}) (RuleType, client.IRule)
type RuleConstructorMap map[string]ConstructorFunc

var constructorMap = RuleConstructorMap{
	"DateRule": func(i interface{}) (RuleType, client.IRule) {
		tArr := i.(map[string]interface{})["limits"].([]interface{})
		rule := &DateRule{"DateRule", [2]int32{int32(tArr[0].(float64)), int32(tArr[1].(float64))}}
		return AppendRule, rule
	},
}

func RegisterRuleConstructor(set RuleConstructorMap) {
	for name, function := range set {
		constructorMap[name] = function
	}
}

type RuleCollection struct {
	appendingRules []client.IRule
	selectRules    []client.IRule
}

func (r *RuleCollection) AppendRules() []client.IRule {
	return r.appendingRules
}

func (r *RuleCollection) SelectRules() []client.IRule {
	return r.selectRules
}

func (r *RuleCollection) getConstructorFor(s map[string]interface{}) ConstructorFunc {
	name := s["name"].(string)
	return constructorMap[name]
}

func (r *RuleCollection) UnmarshalJSON(b []byte) error {
	var raw []map[string]interface{}
	json.Unmarshal(b, &raw)
	for _, entry := range raw {
		if e := r.getConstructorFor(entry); e != nil {
			switch t, rule := e(entry); t {
			case AppendRule:
				r.appendingRules = append(r.appendingRules, rule)
			case SelectRule:
				r.selectRules = append(r.selectRules, rule)
			}
		} else {
			return errors.UnknownRuleError{}
		}
	}

	return nil
}

type DateRule struct {
	Name   string   `json:"name"`
	Limits [2]int32 `json:"limits"`
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
	Id   string `json:""`
}

func (f FollowingRule) GetName() string {
	return f.Name
}

func (f FollowingRule) Validate(i interface{}) (bool, error) {
	//examined := i.(client.HasDateAttribute).GetCreationDate()
	//if examined > d.Limits[1] {
	//	return false, nil
	//}
	//if examined < d.Limits[0] {
	//	return false, errors.ShouldStopIterationError{}
	//}
	return true, nil
}
