package nosj

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type JSON struct {
	jmap map[string]interface{}
}

func Json(rawjson string) (JSON, error) {
	j := JSON{}
	err := json.Unmarshal([]byte(rawjson), &j.jmap)
	if err != nil {
		return JSON{}, fmt.Errorf("parse error: %s", err.Error())
	}
	return j, nil
}

func (j JSON) HasKey(key string) bool {
	_, ok := j.jmap[key]
	return ok
}

func (j JSON) IsString(key string) bool {
	val, _ := j.jmap[key]
	return reflect.TypeOf(val) == reflect.TypeOf("")
}

func (j JSON) IsList(key string) bool {
	val, _ := j.jmap[key]
	kind := reflect.TypeOf(val).Kind()
	return kind == reflect.TypeOf([]interface{}{}).Kind()
}

func (j JSON) IsNum(key string) bool {
	val, _ := j.jmap[key]
	return reflect.TypeOf(val) == reflect.TypeOf(64.4)
}

func (j JSON) IsOb(key string) bool {
	val, _ := j.jmap[key]
	kind := reflect.TypeOf(val).Kind()
	return kind == reflect.TypeOf(map[string]interface{}{}).Kind()
}

func (j JSON) IsBool(key string) bool {
	val, _ := j.jmap[key]
	return reflect.TypeOf(val) == reflect.TypeOf(false)
}
