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
	StatusCode StatusCode `json:"status_code"`
	Body       string     `json:"body"`
}

func NewResponse(code StatusCode, body string) *Response {
	return &Response{
		StatusCode: code,
		Body: body,
	}
}

func (r *Response) Encode() ([]byte, error) {
	data, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("encoding error: %v", err)
	}
	return data, nil
}

func (r *Response) Decode(data []byte) error {
	cleanData := cleanJSONData(data)
	err := json.Unmarshal(cleanData, r)
	if err != nil {
		return fmt.Errorf("decoding error: %v", err)
	}
	return nil
}

func (r *Response) Type() MsgType {
	return RESPONSE
}
