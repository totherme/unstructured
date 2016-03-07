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

	return json.IsOfType(m.typ), nil
}

func (m JSONTypeMatcher) FailureMessage(actual interface{}) (message string) {
	if reflect.TypeOf(actual) != reflect.TypeOf(nosj.JSON{}) {
		return fmt.Sprintf("expected a JSON object -- got a %s", reflect.TypeOf(actual))
	}

	json := actual.(nosj.JSON)
	for _, t := range []string{nosj.JSONBool,
		nosj.JSONString,
		nosj.JSONNum,
		nosj.JSONList,
		nosj.JSONNull,
		nosj.JSONOb} {
		if json.IsOfType(t) {
			return fmt.Sprintf("expected a JSON %s -- got a JSON %s", m.typ, t)
		}
	}

	return fmt.Sprintf("expected a JSON %s -- got some other crazy kind of JSON", m.typ)
}
func (m JSONTypeMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("got a JSON %s, but expected not to", m.typ)
}
