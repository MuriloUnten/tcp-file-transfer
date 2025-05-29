package protocol

import (
	"bytes"
	"unicode"
)

type MsgType string

const (
	REQUEST MsgType = "request"
	RESPONSE MsgType = "response"
	SERVER_SENT_EVENT MsgType = "sse"
	STREAM MsgType = "stream"
)

/*
	Message is the interface for all types of communication in the system.
	The possible forms of messages are the following:
		1. Request
		2. Response
		3. SSE (Server Sent Event)
		4. Stream
*/
type Message interface {
	Type() MsgType
	Encode() ([]byte, error)
	Decode(data []byte) (error)
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
