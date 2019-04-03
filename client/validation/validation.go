package validation

import (
	"encoding/json"
	"fmt"
	"giveaway/data/errors"
	bjson "giveaway/utils/bson"
	"go.mongodb.org/mongo-driver/bson"
)

type RuleType int

const (
	PreconditionRule RuleType = iota
	AppendingRule
	PostconditionRule
	SelectRule
)

type IRule interface {
	fmt.Stringer
	bson.Marshaler
	bson.Unmarshaler
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
	PreconditionRules  []IRule `json:"precondition_rules" bson:"precondition_rules"`
	AppendingRules     []IRule `json:"appending_rules" bson:"appending_rules"`
	PostconditionRules []IRule `json:"postcondition_rules" bson:"postcondition_rules"`
	SelectRules        []IRule `json:"select_rules" bson:"select_rules"`
}

func (r *RuleCollection) getConstructorFor(s map[string]interface{}) ConstructorFunc {
	name := s["name"].(string)
	return constructorMap[name]
}

func (r *RuleCollection) UnmarshalJSON(b []byte) error {
	var m map[string]interface{}
	var raw = make([]map[string]interface{}, 0)
	json.Unmarshal(b, &m)
	if m != nil {
		for _, v := range m {
			if v != nil {
				for _, t := range v.([]interface{}) {
					raw = append(raw, t.(map[string]interface{}))
				}
			}
		}
	} else {
		json.Unmarshal(b, &raw)
	}
	for _, entry := range raw {
		if e := r.getConstructorFor(entry); e != nil {
			switch t, rule := e(entry); t {
			case AppendingRule:
				r.AppendingRules = append(r.AppendingRules, rule)
			case PreconditionRule:
				r.PreconditionRules = append(r.PreconditionRules, rule)
			case PostconditionRule:
				r.PostconditionRules = append(r.PostconditionRules, rule)
			case SelectRule:
				r.SelectRules = append(r.SelectRules, rule)
			}
		} else {
			return errors.UnknownRuleError{}
		}
	}

	return nil
}

func (r *RuleCollection) UnmarshalBSON(b []byte) error {
	return bjson.BSONToStruct(b, r)
}
