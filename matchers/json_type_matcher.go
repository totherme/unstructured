package matchers

import (
	"fmt"
	"reflect"

	"github.com/totherme/nosj"
)

type JSONTypeMatcher struct {
	typ string
}

func BeAnObject() JSONTypeMatcher {
	return JSONTypeMatcher{
		typ: nosj.JSONOb,
	}
}

func BeAString() JSONTypeMatcher {
	return JSONTypeMatcher{
		typ: nosj.JSONString,
	}
}

func (m JSONTypeMatcher) Match(actual interface{}) (success bool, err error) {
	if reflect.TypeOf(actual) != reflect.TypeOf(nosj.JSON{}) {
		return false, fmt.Errorf("actual is not a JSON -- actually of type %s", reflect.TypeOf(actual))
	}
	json := actual.(nosj.JSON)

	switch m.typ {
	case nosj.JSONString:
		return json.IsString(), nil
	case nosj.JSONNum:
		return json.IsNum(), nil
	case nosj.JSONOb:
		return json.IsOb(), nil
	case nosj.JSONList:
		return json.IsList(), nil
	case nosj.JSONNull:
		return json.IsNull(), nil
	case nosj.JSONBool:
		return json.IsBool(), nil
	}
	return false, fmt.Errorf("this is some kind of crazy JSON")
}
func (m JSONTypeMatcher) FailureMessage(actual interface{}) (message string) {
	if reflect.TypeOf(actual) != reflect.TypeOf(nosj.JSON{}) {
		return fmt.Sprintf("expected a JSON object -- got a %s", reflect.TypeOf(actual))
	}

	json := actual.(nosj.JSON)
	if json.IsBool() {
		return fmt.Sprintf("expected a JSON %s -- got a JSON bool", m.typ)
	}
	if json.IsString() {
		return fmt.Sprintf("expected a JSON %s -- got a JSON string", m.typ)
	}
	if json.IsNum() {
		return fmt.Sprintf("expected a JSON %s -- got a JSON number", m.typ)
	}
	if json.IsList() {
		return fmt.Sprintf("expected a JSON %s -- got a JSON list", m.typ)
	}
	if json.IsNull() {
		return fmt.Sprintf("expected a JSON %s -- got JSON null", m.typ)
	}
	if json.IsOb() {
		return fmt.Sprintf("expected a JSON %s -- got JSON object", m.typ)
	}

	return fmt.Sprintf("expected a JSON %s -- got some other crazy kind of JSON", m.typ)
}
func (m JSONTypeMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("got a JSON %s, but expected not to", m.typ)
}
