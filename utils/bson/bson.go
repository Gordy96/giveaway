package bson

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
)

func StructToBSON(i interface{}) ([]byte, error) {
	values := map[string]interface{}{}
	tBts, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(tBts, &values)
	if err != nil {
		return nil, err
	}
	return bson.Marshal(values)
}

func BSONToStruct(b []byte, i interface{}) error {
	values := map[string]interface{}{}
	err := bson.Unmarshal(b, &values)
	if err != nil {
		return err
	}
	tBts, err := json.Marshal(values)
	if err != nil {
		return err
	}
	err = json.Unmarshal(tBts, i)
	if err != nil {
		return err
	}
	return nil
}
