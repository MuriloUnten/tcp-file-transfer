package protocol

import (
	"encoding/json"
	"fmt"
	"bytes"
	"unicode"
)

type Request struct {
	Method Method `json:"method"`
	Body   string `json:"body"`
}

type Response struct {
	StatusCode StatusCode `json:"status_code"`
	Body       string     `json:"body"`
}

func cleanJSONData(data []byte) []byte {
	// Remove null bytes and other control characters
	return bytes.Map(func(r rune) rune {
		if r == 0 || unicode.IsControl(r) {
			return -1
		}
		return r
	}, data)
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
