package protocol

import (
	"encoding/json"
	"fmt"
)

type Stream struct {
	ByteCount int `json:"byte_count"`
	Body   string `json:"body"`
}

func NewStream(count int, body string) *Stream {
	return &Stream{
		ByteCount: count,
		Body: body,
	}
}

func (s *Stream) Encode() ([]byte, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("encoding error: %v", err)
	}
	return data, nil
}

func (s *Stream) Decode(data []byte) error {
	cleanData := cleanJSONData(data)
	err := json.Unmarshal(cleanData, s)
	if err != nil {
		return fmt.Errorf("decoding error: %v", err)
	}
	return nil
}

func (s *Stream) Type() MsgType {
	return STREAM
}
