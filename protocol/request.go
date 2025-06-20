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
	Method Method `json:"method"`
	Body   string `json:"body"`
}

func NewRequest(method Method, body string) *Request {
	return &Request{
		Method: method,
		Body: body,
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

func (m *Request) Type() MsgType {
	return REQUEST
}
