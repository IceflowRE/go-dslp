package v2_0

import (
	"bytes"
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
	msg.Body = append([]byte(time.Now().Format("2006-01-02T15:04:05+07:00")), LineBreak...)
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

func NewGroupNotify(group string, body []string) *Message {
	msg := NewMessage()
	msg.Type = TGroupNotify
	msg.Header = []string{group, strconv.Itoa(len(body))}
	var buf bytes.Buffer
	for _, line := range body {
		buf.Write([]byte(line))
		buf.Write(LineBreak)
	}
	msg.Body = buf.Bytes()
	return msg
}

func NewUserJoin(name string) *Message {
	msg := NewMessage()
	msg.Type = TUserJoin
	msg.Header = []string{name}
	return msg
}

func NewUserLeave(name string) *Message {
	msg := NewMessage()
	msg.Type = TUserLeave
	msg.Header = []string{name}
	return msg
}

func NewUserTextNotify(sender string, target string, body []string) *Message {
	msg := NewMessage()
	msg.Type = TUserTextNotify
	msg.Header = []string{sender, target, strconv.Itoa(len(body))}
	var buf bytes.Buffer
	for _, line := range body {
		buf.Write([]byte(line))
		buf.Write(LineBreak)
	}
	msg.Body = buf.Bytes()
	return msg
}

func NewUserFileNotify(sender string, target string, filename string, mime string, body []byte) *Message {
	msg := NewMessage()
	msg.Type = TUserFileNotify
	msg.Header = []string{sender, target, filename, mime, strconv.Itoa(len(body))}
	msg.Body = body
	return msg
}

func NewErrorMsg(body string) *Message {
	msg := NewMessage()
	msg.Type = TError
	msg.Header = []string{strconv.Itoa(len(body))}
	msg.Body = []byte(body)
	msg.Body = append(msg.Body, LineBreak...)
	return msg
}
