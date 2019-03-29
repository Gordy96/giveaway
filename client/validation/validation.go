package validation

import (
	"encoding/json"
	"giveaway/client/validation/rules"
	"giveaway/data/errors"
)

type RuleType int

const (
	AppendRule RuleType = iota
	SelectRule RuleType = AppendRule + 1
)

type IRule interface {
	Validate(interface{}) (bool, error)
}

type ConstructorFunc func(interface{}) (RuleType, IRule)
type RuleConstructorMap map[string]ConstructorFunc

var constructorMap = RuleConstructorMap{
	"DateRule": func(i interface{}) (RuleType, IRule) {
		tArr := i.(map[string]interface{})["limits"].([]interface{})
		rule := &rules.DateRule{"DateRule", [2]int64{int64(tArr[0].(float64)), int64(tArr[1].(float64))}}
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
