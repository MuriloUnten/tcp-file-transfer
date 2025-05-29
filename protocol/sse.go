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

func (m *SSE) Encode() ([]byte, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("encoding error: %v", err)
	}
	return data, nil
}

func (m *SSE) Decode(data []byte) error {
	cleanData := cleanJSONData(data)
	err := json.Unmarshal(cleanData, m)
	if err != nil {
		return fmt.Errorf("decoding error: %v", err)
	}
	return nil
}
