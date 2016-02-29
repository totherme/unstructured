package nosj

import (
	"encoding/json"
	"fmt"
)

type JSON struct {
	st interface{}
}

func Json(rawjson string) (JSON, error) {
	j := JSON{}
	err := json.Unmarshal([]byte(rawjson), &j.st)
	if err != nil {
		return JSON{}, fmt.Errorf("parse error: %s", err.Error())
	}
	return j, nil
}

func (j JSON) HasKey(key string) bool {
	jmap := j.st.(map[string]interface{})
	_, ok := jmap[key]
	return ok
}
