package messagev1_2

import (
	"net"
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
	content := time.Now().Format("2006-01-02T15:04:05+07:00")
	msg.Content = []byte(content)
	return msg
}

func NewGroupJoin(group string) *Message {
	msg := NewMessage()
	msg.Type = TGroupJoin
	msg.Content = []byte(group)
	return msg
}

func NewGroupLeave(group string) *Message {
	msg := NewMessage()
	msg.Type = TGroupLeave
	msg.Content = []byte(group)
	return msg
}

func NewGroupNotify(group string, content string) *Message {
	msg := NewMessage()
	msg.Type = TGroupNotify
	tmp := group + "\r\n" + content
	msg.Content = []byte(tmp)
	return msg
}

func NewPeerNotfiy(ip net.IP, content string) *Message {
	msg := NewMessage()
	msg.Type = TPeerNotify
	tmp := ip.String() + "\r\n" + content
	msg.Content = []byte(tmp)
	return msg
}

func NewErrorMsg(content string) *Message {
	msg := NewMessage()
	msg.Type = TError
	msg.Content = []byte(content)
	return msg
}
