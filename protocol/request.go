package protocol

import (
	"encoding/json"
	"fmt"
)

type Method string

const (
	Chat Method = "CHAT"
	Fetch Method = "FETCH"
	Quit Method = "QUIT"
)

type Request struct {
	MessageType MsgType `json:"type"`
	Method Method `json:"method"`
	Body   string `json:"body"`
}

func NewRequest() *Request {
	return &Request{
		MessageType: REQUEST,
	}
}

func (m *Request) Encode() ([]byte, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("encoding error: %v", err)
	}
	return data, nil
}

func (m *Request) Decode(data []byte) error {
	cleanData := cleanJSONData(data)
	err := json.Unmarshal(cleanData, m)
	if err != nil {
		return fmt.Errorf("decoding error: %v", err)
	}
	return nil
}
