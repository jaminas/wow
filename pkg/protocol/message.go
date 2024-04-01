package protocol

import (
	"fmt"
	"strconv"
	"strings"
)

type Message struct {
	Header  int
	Payload string
}

// ToString
func (m *Message) ToString() string {
	return fmt.Sprintf("%d|%s", m.Header, m.Payload)
}

// ParseMessage
func ParseMessage(str string) (*Message, error) {
	str = strings.TrimSpace(str)
	var header int

	//
	parts := strings.Split(str, "|")
	if len(parts) < 1 || len(parts) > 2 { //only 1 or 2 parts allowed
		return nil, fmt.Errorf("message doesn't match protocol")
	}

	//
	header, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("parse header error")
	}
	msg := Message{
		Header: header,
	}

	//
	if len(parts) == 2 {
		msg.Payload = parts[1]
	}
	return &msg, nil
}
