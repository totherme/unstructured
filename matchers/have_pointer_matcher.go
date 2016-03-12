package matchers

import (
	"fmt"
	"github.com/totherme/nosj"
)

type HaveJSONPointerMatcher struct {
	p string
}

func HaveJSONPointer(p string) HaveJSONPointerMatcher {
	return HaveJSONPointerMatcher{p: p}
}

func (m HaveJSONPointerMatcher) Match(actual interface{}) (bool, error) {

	switch t := actual.(type) {
	default:
		return false, fmt.Errorf("not a JSON object. Have you done nosj.ParseJSON(...)?")
	case nosj.JSON:
		return t.HasPointer(m.p)
	}
}

func (m HaveJSONPointerMatcher) FailureMessage(actual interface{}) (message string) {
	actualString := fmt.Sprintf("%+v", actual)
	return fmt.Sprintf("expected '%s' to be a nosj.JSON object with pointer '%s'",
		truncateString(actualString),
		m.p)
}

func (m HaveJSONPointerMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	actualString := fmt.Sprintf("%+v", actual)
	return fmt.Sprintf("expected '%s' not to contain the pointer '%s'",
		truncateString(actualString),
		m.p)
}
