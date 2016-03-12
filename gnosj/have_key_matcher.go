package gnosj

import (
	"fmt"

	"github.com/totherme/nosj"
)

type HaveJSONKeyMatcher struct {
	key string
}

func HaveJSONKey(key string) HaveJSONKeyMatcher {
	return HaveJSONKeyMatcher{key: key}
}

func (m HaveJSONKeyMatcher) Match(actual interface{}) (bool, error) {
	switch j := actual.(type) {
	default:
		return false, fmt.Errorf("not a JSON object. Have you done nosj.ParseJSON(...)?")
	case nosj.JSON:
		return j.HasKey(m.key), nil
	}
}

func (m HaveJSONKeyMatcher) FailureMessage(actual interface{}) (message string) {
	actualString := fmt.Sprintf("%+v", actual)
	return fmt.Sprintf("expected '%s' to be a nosj.JSON object with key '%s'",
		truncateString(actualString),
		m.key)
}

func (m HaveJSONKeyMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	actualString := fmt.Sprintf("%+v", actual)
	return fmt.Sprintf("expected '%s' not to contain the key '%s'",
		truncateString(actualString),
		m.key)
}

func truncateString(s string) (t string) {
	if len(s) > 50 {
		t = fmt.Sprintf("%s...", s[0:50])
	} else {
		t = s
	}
	return
}
