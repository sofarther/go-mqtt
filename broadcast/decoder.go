package broadcast

import (
	"bytes"
	"encoding/json"
)

type Message struct {
	Topic string `json:"topic"`
	Payload string `json:"message"`
}

func Decode(m string) (message Message,  err error) {
	err = json.Unmarshal(bytes.NewBufferString(m).Bytes(), &message)

	return
}
