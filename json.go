package nosj

import (
	"encoding/json"
	"fmt"
	"github.com/xeipuuv/gojsonpointer"
	"reflect"
)

const (
	JSONString = "string"
	JSONNum    = "number"
	JSONOb     = "object"
	JSONList   = "list"
	JSONNull   = "null"
	JSONBool   = "bool"
)

type JSON struct {
	nosj interface{}
}

func ParseJSON(rawjson string) (JSON, error) {
	j := JSON{}
	err := json.Unmarshal([]byte(rawjson), &j.nosj)
	if err != nil {
		return JSON{}, fmt.Errorf("parse error: %s", err.Error())
	}
	return j, nil
}

func (j JSON) IsOb() bool {
	return reflect.TypeOf(j.nosj) == reflect.TypeOf(map[string]interface{}{})
}

func (j JSON) HasKey(key string) bool {
	jmap := j.nosj.(map[string]interface{})
	_, ok := jmap[key]
	return ok
}

func (j JSON) HasPointer(p string) bool {
	pointer, err := gojsonpointer.NewJsonPointer(p)
	if err != nil {
		return false
	}
	_, _, err = pointer.Get(j.nosj)
	return err == nil
}

func (j JSON) GetField(key string) JSON {
	jmap := j.nosj.(map[string]interface{})
	val, ok := jmap[key]
	if !ok {
		panic("getting a non-existing field from a JSON")
	}
	return JSON{nosj: val}
}

func (j JSON) F(key string) JSON {
	return j.GetField(key)
}

func (j JSON) GetByPointer(p string) (nosj JSON, err error) {

	pointer, err := gojsonpointer.NewJsonPointer(p)
	if err != nil {
		return
	}
	json, _, err := pointer.Get(j.nosj)
	nosj = JSON{nosj: json}
	return
}

func (j JSON) IsString() bool {
	return reflect.TypeOf(j.nosj) == reflect.TypeOf("")
}

func (j JSON) StringValue() string {
	return j.nosj.(string)
}

func (j JSON) IsNum() bool {
	return reflect.TypeOf(j.nosj) == reflect.TypeOf(64.4)
}

func (j JSON) NumValue() float64 {
	return j.nosj.(float64)
}

func (j JSON) IsBool() bool {
	return reflect.TypeOf(j.nosj) == reflect.TypeOf(true)
}

func (j JSON) BoolValue() bool {
	return j.nosj.(bool)
}

func (j JSON) IsList() bool {
	if j.nosj == nil {
		return false
	}
	return reflect.TypeOf(j.nosj).Kind() == reflect.TypeOf([]interface{}{}).Kind()
}

func (j JSON) ListValue() (list []JSON) {
	list = []JSON{}
	for _, val := range j.nosj.([]interface{}) {
		list = append(list, JSON{nosj: val})
	}
	return
}

func (j JSON) IsNull() bool {
	return j.nosj == nil
}

func (j JSON) IsOfType(typ string) bool {
	switch typ {
	case JSONOb:
		return j.IsOb()
	case JSONString:
		return j.IsString()
	case JSONList:
		return j.IsList()
	case JSONNum:
		return j.IsNum()
	case JSONBool:
		return j.IsBool()
	case JSONNull:
		return j.IsNull()
	default:
		panic("that's not a JSON type!")
	}
}
