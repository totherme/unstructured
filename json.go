package nosj

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type JSON struct {
	st interface{}
}

func Json(rawjson string) (JSON, error) {
	j := JSON{}
	err := json.Unmarshal([]byte(rawjson), &j.st)
	if err != nil {
		return JSON{}, fmt.Errorf("parse error: %s", err.Error())
	}
	return j, nil
}

func (j JSON) IsOb() bool {
	return reflect.TypeOf(j.st) == reflect.TypeOf(map[string]interface{}{})
}

func (j JSON) HasKey(key string) bool {
	jmap := j.st.(map[string]interface{})
	_, ok := jmap[key]
	return ok
}

func (j JSON) GetField(key string) JSON {
	jmap := j.st.(map[string]interface{})
	val, ok := jmap[key]
	if !ok {
		panic("getting a non-existing field from a JSON")
	}
	return JSON{st: val}
}

func (j JSON) IsString() bool {
	return reflect.TypeOf(j.st) == reflect.TypeOf("")
}

func (j JSON) StringValue() string {
	return j.st.(string)
}

func (j JSON) IsNum() bool {
	return reflect.TypeOf(j.st) == reflect.TypeOf(64.4)
}

func (j JSON) NumValue() float64 {
	return j.st.(float64)
}

func (j JSON) IsBool() bool {
	return reflect.TypeOf(j.st) == reflect.TypeOf(true)
}

func (j JSON) BoolValue() bool {
	return j.st.(bool)
}

func (j JSON) IsList() bool {
	if j.st == nil {
		return false
	}
	return reflect.TypeOf(j.st).Kind() == reflect.TypeOf([]interface{}{}).Kind()
}

func (j JSON) ListValue() (list []JSON) {
	list = []JSON{}
	for _, val := range j.st.([]interface{}) {
		list = append(list, JSON{st: val})
	}
	return
}

func (j JSON) IsNull() bool {
	return j.st == nil
}
