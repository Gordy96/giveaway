package validation

import (
	"encoding/json"
	"giveaway/client/validation/rules"
	"giveaway/data/errors"
)

type RuleType int

const (
	AppendRule RuleType = iota
	SelectRule
)

type IRule interface {
	Validate(interface{}) (bool, error)
}

type ConstructorFunc func(interface{}) (RuleType, IRule)
type RuleConstructorMap map[string]ConstructorFunc

var constructorMap = RuleConstructorMap{
	"DateRule": func(i interface{}) (RuleType, IRule) {
		tArr := i.(map[string]interface{})["limits"].([]interface{})
		limits := [2]int64{}
		if len(tArr) == 1 {
			if t := tArr[0]; t != nil {
				limits[0] = int64(t.(float64))
			}
		} else {
			if t := tArr[0]; t != nil {
				limits[0] = int64(t.(float64))
			}
			if t := tArr[1]; t != nil {
				limits[1] = int64(t.(float64))
			}
		}
		rule := &rules.DateRule{"DateRule", limits}
		return AppendRule, rule
	},
}

func RegisterRuleConstructor(set RuleConstructorMap) {
	for name, function := range set {
		constructorMap[name] = function
	}
}

type RuleCollection struct {
	appendingRules []IRule
	selectRules    []IRule
}

func (r *RuleCollection) AppendRules() []IRule {
	return r.appendingRules
}

func (r *RuleCollection) SelectRules() []IRule {
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
