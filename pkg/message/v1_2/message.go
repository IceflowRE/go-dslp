package messagev1_2

import (
	"bytes"
	"errors"
	"net"
	"regexp"
	"strings"

	"github.com/IceflowRE/go-dslp/pkg/message"
)

var rxMessage = regexp.MustCompile(`(?:^|\r\n)(?:dslp\/1\.2)\r\n(` + strings.Join(Types, "|") + `)\r\n(?:((?:.|\r|\n)*?)\r\n)?(?:dslp\/end)\r\n`)

const (
	TRequestTime  = "request time"
	TResponseTime = "response time"
	TGroupJoin    = "group join"
	TGroupLeave   = "group leave"
	TGroupNotify  = "group notify"
	TPeerNotify   = "peer notify"
	TError        = "error"
)

var Types = []string{TRequestTime, TResponseTime, TGroupJoin, TGroupLeave, TGroupNotify, TPeerNotify, TError}

type Message struct {
	message.IMessage
	Header string
	Type   string
	// excludes the ending \r\n
	Content []byte
	End     string
}

func NewMessage() *Message {
	return &Message{
		Header: "dslp/1.2",
		End:    "dslp/end",
	}
}

func (msg *Message) GetType() string {
	return msg.Type
}

func (msg *Message) GetContent() *string {
	if msg.Content != nil {
		tmp := string(msg.Content)
		return &tmp
	}
	return nil
}

func (msg *Message) ToBytes() []byte {
	var buf bytes.Buffer
	buf.WriteString(msg.Header)
	buf.WriteString("\r\n")
	buf.WriteString(msg.Type)
	buf.WriteString("\r\n")
	if msg.Content != nil {
		buf.Write(msg.Content)
		buf.WriteString("\r\n")
	}
	buf.WriteString(msg.End)
	buf.WriteString("\r\n")
	return buf.Bytes()
}

func (msg *Message) Valid() error {
	var errMsg string
	switch msg.GetType() {
	case TRequestTime:
		if msg.GetContent() != nil {
			errMsg = "requires to have an empty body"
		}
	case TResponseTime:
		// no need to care about validity
		return nil
	case TGroupJoin:
		if msg.Content == nil || len(strings.Split(*msg.GetContent(), "\r\n")) != 1 {
			errMsg = "must have one data line"
		}
	case TGroupLeave:
		if msg.Content == nil || len(strings.Split(*msg.GetContent(), "\r\n")) != 1 {
			errMsg = "must have one data line"
		}
	case TGroupNotify:
		if msg.Content == nil {
			errMsg = "must have at least two data lines"
			break
		}
		split := strings.Split(*msg.GetContent(), "\r\n")
		if len(split) < 2 {
			errMsg = "must have at least two data lines"
		} else if split[0] == "" || split[1] == "" {
			errMsg = "the first two lines cannot be empty"
		}
	case TPeerNotify:
		if msg.Content == nil {
			errMsg = "must have at least two data lines"
			break
		}
		split := strings.SplitN(*msg.GetContent(), "\r\n", 2)
		if len(split) != 2 {
			errMsg = "must have at least two data lines"
		} else if split[0] == "" {
			errMsg = "the first line cannot be empty"
		} else if net.ParseIP(split[0]) == nil {
			errMsg = "IP had a wrong format"
		}
	case TError:
		if msg.Content == nil || len(strings.Split(*msg.GetContent(), "\r\n")) < 1 {
			errMsg = "must have at least one data line"
		}
	default:
		return errors.New("Message type (" + msg.Type + ") is invalid.")
	}
	if errMsg != "" {
		return errors.New("type (" + msg.Type + ") " + errMsg)
	}
	return nil
}
