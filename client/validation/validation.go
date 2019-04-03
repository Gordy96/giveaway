package validation

import (
	"encoding/json"
	"fmt"
	"giveaway/data/errors"
)

type RuleType int

const (
	PreconditionRule RuleType = iota
	AppendRule
	PostconditionRule
	SelectRule
)

type IRule interface {
	fmt.Stringer
	Validate(interface{}) (bool, error)
}

type ConstructorFunc func(interface{}) (RuleType, IRule)
type RuleConstructorMap map[string]ConstructorFunc

var constructorMap = RuleConstructorMap{}

func RegisterRuleConstructor(set RuleConstructorMap) {
	for name, function := range set {
		constructorMap[name] = function
	}
}

type RuleCollection struct {
	preconditionRules  []IRule
	appendingRules     []IRule
	postconditionRules []IRule
	selectRules        []IRule
}

func (r *RuleCollection) AppendRules() []IRule {
	return r.appendingRules
}

func (r *RuleCollection) SelectRules() []IRule {
	return r.selectRules
}

func (r *RuleCollection) PreconditionRules() []IRule {
	return r.preconditionRules
}

func (r *RuleCollection) PostconditionRules() []IRule {
	return r.postconditionRules
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
			case PreconditionRule:
				r.preconditionRules = append(r.preconditionRules, rule)
			case PostconditionRule:
				r.postconditionRules = append(r.postconditionRules, rule)
			case SelectRule:
				r.selectRules = append(r.selectRules, rule)
			}
		} else {
			return errors.UnknownRuleError{}
		}
	}

	return nil
}
