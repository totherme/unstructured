package nosj

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type JSON struct {
	st   interface{}
	jmap map[string]interface{}
}

func Json(rawjson string) (JSON, error) {
	j := JSON{}
	err := json.Unmarshal([]byte(rawjson), &j.st)
	if err != nil {
		return JSON{}, fmt.Errorf("parse error: %s", err.Error())
	}
	return j, nil
}

func (j JSON) HasKey(key string) bool {
	jmap := j.st.(map[string]interface{})
	_, ok := jmap[key]
	return ok
}

func (j JSON) IsString(key string) bool {
	jmap := j.st.(map[string]interface{})
	val, _ := jmap[key]
	return reflect.TypeOf(val) == reflect.TypeOf("")
}

func (j JSON) IsList(key string) bool {
	jmap := j.st.(map[string]interface{})
	val, _ := jmap[key]
	kind := reflect.TypeOf(val).Kind()
	return kind == reflect.TypeOf([]interface{}{}).Kind()
}

func (j JSON) IsNum(key string) bool {
	jmap := j.st.(map[string]interface{})
	val, _ := jmap[key]
	return reflect.TypeOf(val) == reflect.TypeOf(64.4)
}

func (j JSON) IsOb(key string) bool {
	jmap := j.st.(map[string]interface{})
	val, _ := jmap[key]
	kind := reflect.TypeOf(val).Kind()
	return kind == reflect.TypeOf(map[string]interface{}{}).Kind()
}

func (j JSON) IsBool(key string) bool {
	jmap := j.st.(map[string]interface{})
	val, _ := jmap[key]
	return reflect.TypeOf(val) == reflect.TypeOf(false)
}

func (j JSON) GetString(key string) string {
	jmap := j.st.(map[string]interface{})
	val, _ := jmap[key]
	return val.(string)
}

func (j JSON) GetList(key string) []JSON {
	jmap := j.st.(map[string]interface{})
	rawlist := jmap[key].([]interface{})
	jlist := []JSON{}
	for ob := range rawlist {
		jOb := JSON{}
		jOb.st = ob
		jlist = append(jlist, jOb)
	}
	return jlist
}
