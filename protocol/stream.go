package protocol

import (
	"encoding/json"
	"fmt"
)

type Stream struct {
	MessageType MsgType `json:"type"`
	ByteCount int `json:"byte_count"`
	Body   string `json:"body"`
}

func NewStream() *Stream {
	return &Stream{
		MessageType: STREAM,
	}
}

func (m *Stream) Encode() ([]byte, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("encoding error: %v", err)
	}
	return data, nil
}

func (m *Stream) Decode(data []byte) error {
	cleanData := cleanJSONData(data)
	err := json.Unmarshal(cleanData, m)
	if err != nil {
		return fmt.Errorf("decoding error: %v", err)
	}
	return nil
}
