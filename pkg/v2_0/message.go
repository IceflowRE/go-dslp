package v2_0

import (
	"bytes"
	"errors"
	"net"
	"regexp"
	"strconv"
	"strings"

	"github.com/IceflowRE/go-dslp/pkg/message"
	"github.com/IceflowRE/go-dslp/pkg/utils"
)

var RxHeader = regexp.MustCompile(`(?:^|\r\n)(?:dslp\/2\.0)\r\n(` + strings.Join(Types, "|") + `)\r\n(?:((?:.|\r|\n)*?)\r\n)?(?:dslp\/body)\r\n`)
var LineBreak = []byte("\r\n")

const (
	TRequestTime    = "request time"
	TResponseTime   = "response time"
	TGroupJoin      = "group join"
	TGroupLeave     = "group leave"
	TGroupNotify    = "group notify"
	TUserJoin       = "user join"
	TUserLeave      = "user leave"
	TUserTextNotify = "user text notify"
	TUserFileNotify = "user file notify"
	TError          = "error"
)

var Types = []string{TRequestTime, TResponseTime, TGroupJoin, TGroupLeave, TGroupNotify, TUserJoin, TUserLeave, TUserTextNotify, TUserFileNotify, TError}

type Message struct {
	message.IMessage // just as information that it implements that
	HeaderKey string
	BodyKey   string
	Type      string
	// header slice excludes the type
	Header []string
	Body   []byte
}

func NewMessage() *Message {
	return &Message{
		HeaderKey: "dslp/2.0",
		BodyKey:   "dslp/body",
	}
}

func (msg *Message) GetType() string {
	return msg.Type
}

func (msg *Message) GetContent() *string {
	if msg.Body != nil {
		tmp := string(msg.Body)
		return &tmp
	}
	return nil
}

func (msg *Message) GetRawContent() []byte {
	return msg.Body
}

func (msg *Message) ToBytes() []byte {
	msg.UpdateHeader()

	var buf bytes.Buffer
	buf.WriteString(msg.HeaderKey)
	buf.WriteString("\r\n")
	buf.WriteString(msg.Type)
	buf.WriteString("\r\n")
	if msg.Header != nil {
		for _, line := range msg.Header {
			buf.WriteString(line)
			buf.WriteString("\r\n")
		}
	}
	buf.WriteString(msg.BodyKey)
	buf.WriteString("\r\n")
	if msg.Body != nil {
		buf.Write(msg.Body)
		buf.WriteString("\r\n")

	}
	return buf.Bytes()
}

func (msg *Message) UpdateHeader() error {
	switch msg.Type {
	case TRequestTime:
		msg.Header = nil
	case TResponseTime:
		msg.Header = nil
	case TGroupJoin: // cannot create header for this
		return nil
	case TGroupLeave: // cannot create header for this
		return nil
	case TGroupNotify:
		if msg.Header == nil || len(msg.Header) != 2 {
			newHeader := make([]string, 2)
			if len(msg.Header) > 0 {
				newHeader[0] = msg.Header[0]
			}
			newHeader[1] = strconv.Itoa(len(msg.Body))
			msg.Header = newHeader
		}
	//case TUserJoin:
	//case TUserLeave:
	//case TUserTextNotify:
	//case TUserFileNotify:
	case TError:
		newHeader := make([]string, 1)
		newHeader[0] = strconv.Itoa(len(msg.Body))
		msg.Header = newHeader
	default:
		errors.New("don't know how to make a header for type " + msg.Type)
	}
	return nil
}

func (msg *Message) Valid() error {
	var errMsg string
	switch msg.Type {
	case TRequestTime:
		if msg.Header != nil {
			errMsg = "header must have no additional data"
		} else if msg.Body != nil {
			errMsg = "must have an empty body"
		}
	case TResponseTime:
		if msg.Header != nil {
			errMsg = "header must have no additional data"
		} else if msg.Body == nil || len(msg.Body) != 1 {
			errMsg = "must have one body line"
		}
	case TGroupJoin:
		if msg.Header == nil || len(msg.Header) != 1 || msg.Header[0] == "" {
			errMsg = "header must contain the group to join"
		} else if msg.Body != nil {
			errMsg = "must have an empty body"
		}
	case TGroupLeave:
		if msg.Header == nil || len(msg.Header) != 1 || msg.Header[0] == "" {
			errMsg = "header must contain the group to leave"
		} else if msg.Body != nil {
			errMsg = "must have an empty body"
		}
	case TGroupNotify:
		if msg.Header == nil || len(msg.Header) != 2 || msg.Header[0] == "" {
			errMsg = "header must contain the group to notify and the correct number of lines"
		} else if lines, err := strconv.Atoi(msg.Header[1]); err != nil {
			errMsg = "header must contain the group to notify and the correct number of lines"
		} else if lines != len(msg.Body) {
			errMsg = "body does not match the size written in header"
		}
	case TError:
		if msg.Header == nil || len(msg.Header) != 1 {
			errMsg = "header must contain the correct number of lines"
		} else if lines, err := strconv.Atoi(msg.Header[0]); err != nil {
			errMsg = "header must contain the correct number of lines"
		} else if lines != len(msg.Body) {
			errMsg = "body does not match the size written in header"
		}
	default:
		return errors.New("Message type (" + msg.Type + ") is invalid.")
	}
	if errMsg != "" {
		return errors.New("type (" + msg.Type + ") " + errMsg)
	}
	return nil
}

// HandleMessage requires a valid message
func HandleMessage(msg message.IMessage, conn net.Conn) error {
	switch msg.GetType() {
	case TRequestTime:
		message.SendMessage(NewResponseTimeMsg(), conn)
	case TResponseTime:
		// do nothing
	case TGroupJoin:
		JoinGroup(conn, *msg.GetContent())
	case TGroupLeave:
		return LeaveGroup(conn, *msg.GetContent())
	case TGroupNotify:
		split := strings.SplitN(*msg.GetContent(), "\r\n", 2)
		SendToGroup(split[0], split[1])
	case TError:
		utils.Println(conn, "Error message received", *msg.GetContent())
	default:
		utils.Println(conn, "type invalid", msg.GetType())
	}
	return nil
}
