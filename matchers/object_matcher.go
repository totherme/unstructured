package matchers

import (
	"fmt"
	"reflect"

	"github.com/totherme/nosj"
)

type ObjectMatcher struct{}

func BeAnObject() ObjectMatcher {
	return ObjectMatcher{}
}

func (m ObjectMatcher) Match(actual interface{}) (success bool, err error) {
	if reflect.TypeOf(actual) != reflect.TypeOf(nosj.JSON{}) {
		return false, fmt.Errorf("actual is not a JSON -- actually of type %s", reflect.TypeOf(actual))
	}

	json := actual.(nosj.JSON)
	return json.IsOb(), nil
}
func (m ObjectMatcher) FailureMessage(actual interface{}) (message string) {
	if reflect.TypeOf(actual) != reflect.TypeOf(nosj.JSON{}) {
		return fmt.Sprintf("expected a JSON object -- got a %s", reflect.TypeOf(actual))
	}

	json := actual.(nosj.JSON)
	if json.IsBool() {
		return "expected a JSON object -- got a JSON bool"
	}
	if json.IsString() {
		return "expected a JSON object -- got a JSON string"
	}
	if json.IsNum() {
		return "expected a JSON object -- got a JSON number"
	}
	if json.IsList() {
		return "expected a JSON object -- got a JSON list"
	}
	if json.IsNull() {
		return "expected a JSON object -- got JSON null"
	}

	return "expected a JSON object -- got some other crazy kind of JSON"
}
func (m ObjectMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return "got a JSON object, but expected not to"
}
