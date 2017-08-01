// Package nosj provides lightweight JSON reflection for testing.
package nosj

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/ghodss/yaml"
	"github.com/xeipuuv/gojsonpointer"
)

const (
	JSONString = "string"
	JSONNum    = "number"
	JSONOb     = "object"
	JSONList   = "list"
	JSONNull   = "null"
	JSONBool   = "bool"
)

// JSON represents a parsed json structure.
type JSON struct {
	nosj interface{}
}

// ParseYAML unmarshals json from an input string. Use this for generating a YAML struct, whose contents you can examine
// using the following functions.
func ParseYAML(rawjson string) (JSON, error) {
	jsonbytes, err := yaml.YAMLToJSON([]byte(rawjson))
	if err != nil {
		return JSON{}, err
	}
	return ParseJSON(string(jsonbytes))
}

// ParseJSON unmarshals yaml from an input string. Use this for generating a JSON struct, whose contents you can examine
// using the following functions.
func ParseJSON(rawjson string) (JSON, error) {
	j := JSON{}
	err := json.Unmarshal([]byte(rawjson), &j.nosj)
	if err != nil {
		return JSON{}, fmt.Errorf("parse error: %s", err.Error())
	}
	return j, nil
}

// IsOb returns true iff the json represented by this JSON struct is a json object.
func (j JSON) IsOb() bool {
	return reflect.TypeOf(j.nosj) == reflect.TypeOf(map[string]interface{}{})
}

// ObValue returns a golang map[string]interface{} represenation of the json object represented by this JSON struct. If the JSON
// struct does not represent a json object, this method panics. If in doubt, check with `IsOb()`
func (j JSON) ObValue() map[string]interface{} {
	return j.nosj.(map[string]interface{})
}

// HasKey returns true iff the json object represented by this JSON struct contains `key`
//
// Note: this will panic if the json represented by this JSON struct is not a json object. If in doubt, check with `IsOb()`
func (j JSON) HasKey(key string) bool {
	jmap := j.nosj.(map[string]interface{})
	_, ok := jmap[key]
	return ok
}

// HasPointer returns true iff the json object represented by this JSON struct contains the pointer `p`
//
// For more information on json pointers, see https://tools.ietf.org/html/rfc6901
func (j JSON) HasPointer(p string) (bool, error) {
	pointer, err := gojsonpointer.NewJsonPointer(p)
	if err != nil {
		return false, err
	}
	_, _, err = pointer.Get(j.nosj)
	return err == nil, nil
}

// GetField returns a JSON struct containing the contents of the original json at the given `key`. If this method name
// feels too long, use `F(key)`.
//
// Note: this function panics if the given `key` does not exist. If in doubt, check with `HasKey()`.
func (j JSON) GetField(key string) JSON {
	jmap := j.nosj.(map[string]interface{})
	val, ok := jmap[key]
	if !ok {
		panic("getting a non-existing field from a JSON")
	}
	return JSON{nosj: val}
}

// F is a shorthand for `GetField`
func (j JSON) F(key string) JSON {
	return j.GetField(key)
}

// SetField updates the field `fieldName` of this JSON object. If this is not a JSON object, we might crash.
// If the field `fieldName` does not exist on this object, create it.
func (j JSON) SetField(fieldName string, val interface{}) {
	jmap := j.nosj.(map[string]interface{})
	jmap[fieldName] = val
}

// GetByPointer returns a JSON struct containing the contents of the original json at the given pointer address `p`.
// For more information on json pointers, see https://tools.ietf.org/html/rfc6901
func (j JSON) GetByPointer(p string) (nosj JSON, err error) {

	pointer, err := gojsonpointer.NewJsonPointer(p)
	if err != nil {
		return
	}
	json, _, err := pointer.Get(j.nosj)
	nosj = JSON{nosj: json}
	return
}

// RawValue returns the raw go value of the parsed json, without any type checking
func (j JSON) RawValue() interface{} {
	return j.nosj
}

// IsString returns true iff the json represented by this JSON struct is a string.
func (j JSON) IsString() bool {
	return reflect.TypeOf(j.nosj) == reflect.TypeOf("")
}

// StringValue returns the golang string representation of the json string represented by this JSON struct. If the JSON
// struct does not represent a json string, this method panics. If in doubt, check with `IsString()`
func (j JSON) StringValue() string {
	return j.nosj.(string)
}

// IsNum returns true iff the json represented by this JSON struct is a number.
func (j JSON) IsNum() bool {
	return reflect.TypeOf(j.nosj) == reflect.TypeOf(64.4)
}

// NumValue returns the golang float64 representation of the json number represented by this JSON struct. If the JSON
// struct does not represent a json number, this method panics. If in doubt, check with `IsNum()`
func (j JSON) NumValue() float64 {
	return j.nosj.(float64)
}

// IsBool returns true iff the json represented by this JSON struct is a boolean.
func (j JSON) IsBool() bool {
	return reflect.TypeOf(j.nosj) == reflect.TypeOf(true)
}

// BoolValue returns the golang bool representation of the json bool represented by this JSON struct. If the JSON
// struct does not represent a json bool, this method panics. If in doubt, check with `IsBool()`
func (j JSON) BoolValue() bool {
	return j.nosj.(bool)
}

// IsList returns true iff the json represented by this JSON struct is a json list.
func (j JSON) IsList() bool {
	if j.nosj == nil {
		return false
	}
	return reflect.TypeOf(j.nosj).Kind() == reflect.TypeOf([]interface{}{}).Kind()
}

// ListValue returns a golang slice of JSON structs representing the json list represented by this JSON struct.
// If the JSON struct does not represent a json list, this method panics. If in doubt, check with `IsList()`
func (j JSON) ListValue() (list []JSON) {
	list = []JSON{}
	for _, val := range j.nosj.([]interface{}) {
		list = append(list, JSON{nosj: val})
	}
	return
}

// SetElem sets the element at a given index in this JSON list to the given value.
// If this JSON object does not represent a list, return an error
func (j JSON) SetElem(index int, value interface{}) error {
	if !j.IsList() {
		return fmt.Errorf("This is not a list, so you can't set an element of it")
	}
	j.nosj.([]interface{})[index] = value
	return nil
}

// IsNull returns true iff the json represented by this JSON struct is json null.
func (j JSON) IsNull() bool {
	return j.nosj == nil
}

// IsOfType returns true iff the JSON struct represents a json of type `typ`. Valid values of `typ` are listed as constants above.
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
