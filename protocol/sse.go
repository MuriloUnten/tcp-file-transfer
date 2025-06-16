package protocol

import (
	"encoding/json"
	"fmt"
)

type SSE struct {
	MessageType MsgType `json:"type"`
	Body   string `json:"body"`
}

func NewSSE() *SSE {
	return &SSE{
		MessageType: SERVER_SENT_EVENT,
	}
}

func (sse *SSE) Encode() ([]byte, error) {
	data, err := json.Marshal(sse)
	if err != nil {
		return nil, fmt.Errorf("encoding error: %v", err)
	}
	return data, nil
}

func (sse *SSE) Decode(data []byte) error {
	cleanData := cleanJSONData(data)
	err := json.Unmarshal(cleanData, sse)
	if err != nil {
		return fmt.Errorf("decoding error: %v", err)
	}
	return nil
}

func (sse *SSE) Type() MsgType {
	return  SERVER_SENT_EVENT
}
