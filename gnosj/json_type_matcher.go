// Package gnosj provides gomega matchers for using nosj with the ginkgo and
// gomega testing libraries:
//
// https://github.com/onsi/gomega
// https://github.com/onsi/ginkgo
package gnosj

import (
	"fmt"
	"reflect"

	"github.com/totherme/nosj"
)

// JSONTypeMatcher is a gomega matcher which tests if a given value represents
// json data of a given type.
type JSONTypeMatcher struct {
	typ string
}

// BeAnObject returns a gomega matcher which tests if a given value represents
// a json object.
func BeAnObject() JSONTypeMatcher {
	return JSONTypeMatcher{
		typ: nosj.JSONOb,
	}
}

// BeAnObject returns a gomega matcher which tests if a given value represents
// a json string.
func BeAString() JSONTypeMatcher {
	return JSONTypeMatcher{
		typ: nosj.JSONString,
	}
}

// BeAList returns a gomega matcher which tests if a given value represents
// a json list.
func BeAList() JSONTypeMatcher {
	return JSONTypeMatcher{
		typ: nosj.JSONList,
	}
}

// BeANum returns a gomega matcher which tests if a given value represents
// a json num.
func BeANum() JSONTypeMatcher {
	return JSONTypeMatcher{
		typ: nosj.JSONNum,
	}
}

// BeABool returns a gomega matcher which tests if a given value represents
// a json bool.
func BeABool() JSONTypeMatcher {
	return JSONTypeMatcher{
		typ: nosj.JSONBool,
	}
}

// BeNull returns a gomega matcher which tests if a given value represents
// json null.
func BeANull() JSONTypeMatcher {
	return JSONTypeMatcher{
		typ: nosj.JSONNull,
	}
}

// Match is the gomega function that actually checks if the given value is of
// the appropriate json type.
func (m JSONTypeMatcher) Match(actual interface{}) (success bool, err error) {
	switch json := actual.(type) {
	default:
		return false, fmt.Errorf("actual is not a JSON -- actually of type %s", reflect.TypeOf(actual))
	case nosj.JSON:
		return json.IsOfType(m.typ), nil
	}
}

// FailureMessage constructs a hopefully-helpful error message in the case that
// the given value is not of the appropriate json type.
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

// NegatedFailureMessage constructs a hopefully-helpful error message in the
// case that the given value is unexpectedly of the appropriate json type.
func (m JSONTypeMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("got a JSON %s, but expected not to", m.typ)
}
