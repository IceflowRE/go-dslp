package v2_0

import (
	"strconv"
	"time"
)

func NewRequestTime() *Message {
	msg := NewMessage()
	msg.Type = TRequestTime
	return msg
}

func NewResponseTimeMsg() *Message {
	msg := NewMessage()
	msg.Type = TResponseTime
	msg.Body = []byte(time.Now().Format("2006-01-02T15:04:05+07:00"))
	return msg
}

func NewGroupJoin(group string) *Message {
	msg := NewMessage()
	msg.Type = TGroupJoin
	msg.Header = []string{group}
	return msg
}

func NewGroupLeave(group string) *Message {
	msg := NewMessage()
	msg.Type = TGroupLeave
	msg.Header = []string{group}
	return msg
}

func NewGroupNotify(group string, body string) *Message {
	msg := NewMessage()
	msg.Type = TGroupNotify
	msg.Header = []string{group, strconv.Itoa(len(body))}
	msg.Body = []byte(body)
	return msg
}

func NewErrorMsg(body string) *Message {
	msg := NewMessage()
	msg.Type = TError
	msg.Header = []string{strconv.Itoa(len(body))}
	msg.Body = []byte(body)
	return msg
}
