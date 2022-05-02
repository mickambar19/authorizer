package model

import (
	"encoding/json"
)

type Event struct {
	Type string
	Data map[string]interface{}
}

func (e *Event) UnmarshalJSON(b []byte) error {
	var temp interface{}
	json.Unmarshal(b, &temp)
	eventMap := temp.(map[string]interface{})
	for key, v := range eventMap {
		e.Type = key
		e.Data = v.(map[string]interface{})
	}
	return nil
}
