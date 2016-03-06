package matchers

import (
	"fmt"
	"reflect"

	"github.com/totherme/nosj"
)

type StringMatcher struct{}

func BeAString() StringMatcher {
	return StringMatcher{}
}

func (m StringMatcher) Match(actual interface{}) (success bool, err error) {
	if reflect.TypeOf(actual) != reflect.TypeOf(nosj.JSON{}) {
		return false, fmt.Errorf("actual is not a JSON -- actually of type %s", reflect.TypeOf(actual))
	}

	json := actual.(nosj.JSON)
	return json.IsString(), nil
}
func (m StringMatcher) FailureMessage(actual interface{}) (message string) {
	if reflect.TypeOf(actual) != reflect.TypeOf(nosj.JSON{}) {
		return fmt.Sprintf("expected a JSON object -- got a %s", reflect.TypeOf(actual))
	}

	json := actual.(nosj.JSON)
	if json.IsBool() {
		return "expected a JSON string -- got a JSON bool"
	}
	if json.IsOb() {
		return "expected a JSON string -- got a JSON object"
	}
	if json.IsNum() {
		return "expected a JSON string -- got a JSON number"
	}
	if json.IsList() {
		return "expected a JSON string -- got a JSON list"
	}
	if json.IsNull() {
		return "expected a JSON string -- got JSON null"
	}

	return "expected a JSON string -- got some other crazy kind of JSON"
}
func (m StringMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return "got a JSON string, but expected not to"
}
