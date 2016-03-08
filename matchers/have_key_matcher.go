package matchers

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
