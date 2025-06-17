package protocol

import (
	"bytes"
	"errors"
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

func DecodeMessage(b []byte) (Message, error) {
	before, after, found := bytes.Cut(b, []byte("|"))
	if !found {
		return nil, errors.New("invalid message format") // TODO fix this (create defined error types)
	}

	msgType := string(before)
	var msg Message

	switch msgType {
	case string(REQUEST):
		msg = new(Request)
	case string(RESPONSE):
		msg = new(Response)
	case string(SERVER_SENT_EVENT):
		msg = new(SSE)
	case string(STREAM):
		msg = new(SSE)
	default:
		return nil, errors.New("invalid message type " + msgType) // TODO fix this (create defined error types)
	}

	err := msg.Decode(after)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func EncodeMessage(m Message) ([]byte, error) {
	msgBytes := []byte{}
	msgType := []byte{}
	switch m.(type) {
	case *Request:
		msgType = []byte(string(REQUEST))
	case *Response:
		msgType = []byte(string(RESPONSE))
	case *SSE:
		msgType = []byte(string(SERVER_SENT_EVENT))
	case *Stream:
		msgType = []byte(string(STREAM))
	default:
		return nil, errors.New("invalid message type") // TODO fix this (create defined error types)
	}

	msgBytes = append(msgBytes, msgType...)
	msgBytes = append(msgBytes, '\n')

	msgContent, err := m.Encode()
	if err != nil {
		return nil, err
	}

	msgBytes = append(msgBytes, msgContent...)
	return msgBytes, nil
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
