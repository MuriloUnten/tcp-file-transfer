package protocol

import (
	"encoding/json"
	"fmt"
)

type StatusCode string

const (
	Ok StatusCode = "OK"
	BadRequest StatusCode = "BAD_REQUEST"
	NotFound StatusCode = "NOT_FOUND"
	InternalError StatusCode = "INTERNAL_ERROR"
)

type Response struct {
	MessageType MsgType `json:"type"`
	StatusCode StatusCode `json:"status_code"`
	Body       string     `json:"body"`
}

func NewResponse() *Response {
	return &Response{
		MessageType: RESPONSE,
	}
}

func (m *Response) Encode() ([]byte, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("encoding error: %v", err)
	}
	return data, nil
}

func (m *Response) Decode(data []byte) error {
	cleanData := cleanJSONData(data)
	err := json.Unmarshal(cleanData, m)
	if err != nil {
		return fmt.Errorf("decoding error: %v", err)
	}
	return nil
}
