package event

import (
	"encoding/json"
)

func Serialize(evt *Event) ([]byte, error) {
	return json.Marshal(evt)
}

func Deserialize(data []byte) (*Event, error) {
	var evt Event
	err := json.Unmarshal(data, &evt)
	return &evt, err
} 