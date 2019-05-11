package serverv2_0

import (
	"bytes"
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/IceflowRE/go-dslp/pkg/message"
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
	// includes the probably last line ending \r\n in difference to the 1.2
	Body []byte
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
			msg.Header = newHeader
		}
		msg.Header[1] = strconv.Itoa(bytes.Count(msg.Body, LineBreak))
	case TUserJoin: // cannot create header for this
		return nil
	case TUserLeave: // cannot create header for this
		return nil
	case TUserTextNotify:
		if msg.Header == nil || len(msg.Header) != 3 {
			newHeader := make([]string, 3)
			for idx, val := range msg.Header {
				newHeader[idx] = val
				if idx == 1 {
					break
				}
			}
			msg.Header = newHeader
		}
		msg.Header[2] = strconv.Itoa(bytes.Count(msg.Body, LineBreak))
	case TUserFileNotify:
		if msg.Header == nil || len(msg.Header) != 5 {
			newHeader := make([]string, 5)
			for idx, val := range msg.Header {
				newHeader[idx] = val
				if idx == 3 {
					break
				}
			}
			msg.Header = newHeader
		}
		msg.Header[4] = strconv.Itoa(len(msg.Body))
	case TError:
		msg.Header = make([]string, 1)
		msg.Header[0] = strconv.Itoa(len(msg.Body))
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
		// do not check time format
	case TResponseTime:
		if msg.Header != nil {
			errMsg = "header must have no additional data"
		} else if msg.Body == nil || bytes.Count(msg.Body, LineBreak) != 1 {
			errMsg = "must have one body line"
		}
	case TGroupJoin:
		if msg.Header == nil || len(msg.Header) != 1 {
			errMsg = "header does not contain all required data"
		} else if msg.Header[0] == "" {
			errMsg = "group name cannot be empty"
		} else if strings.HasPrefix(msg.Header[0], "dslp/") {
			errMsg = "group name cannot begin with 'dslp/'"
		} else if msg.Body != nil {
			errMsg = "must have an empty body"
		}
	case TGroupLeave:
		if msg.Header == nil || len(msg.Header) != 1 {
			errMsg = "header does not contain all required data"
		} else if msg.Header[0] == "" {
			errMsg = "group name cannot be empty"
		} else if msg.Body != nil {
			errMsg = "must have an empty body"
		}
	case TGroupNotify:
		if msg.Header == nil || len(msg.Header) != 2 {
			errMsg = "header does not contain all required data"
		} else if msg.Header[0] == "" {
			errMsg = "group name cannot be empty"
		} else if lines, err := strconv.Atoi(msg.Header[1]); err != nil {
			errMsg = "header must contain the body size"
		} else if bytes.Count(msg.Body, LineBreak) != lines {
			errMsg = "body does not match the size written in header"
		}
	case TUserJoin:
		if msg.Header == nil || len(msg.Header) != 1 || msg.Header[0] == "" {
			errMsg = "header does not contain all required data"
		} else if msg.Header[0] == "" {
			errMsg = "username cannot be empty"
		} else if strings.HasPrefix(msg.Header[0], "dslp/") {
			errMsg = "username cannot begin with 'dslp/'"
		} else if msg.Body != nil {
			errMsg = "must have an empty body"
		}
	case TUserLeave:
		if msg.Header == nil || len(msg.Header) != 1 || msg.Header[0] == "" {
			errMsg = "header does not contain all required data"
		} else if msg.Header[0] == "" {
			errMsg = "username cannot be empty"
		} else if msg.Body != nil {
			errMsg = "must have an empty body"
		}
	case TUserTextNotify:
		if msg.Header == nil || len(msg.Header) != 3 {
			errMsg = "header does not contain all required data"
		} else if msg.Header[0] == "" {
			errMsg = "sender name cannot be empty"
		} else if msg.Header[1] == "" {
			errMsg = "target name cannot be empty"
		} else if lines, err := strconv.Atoi(msg.Header[2]); err != nil {
			errMsg = "header must contain the body size"
		} else if bytes.Count(msg.Body, LineBreak) != lines {
			errMsg = "body does not match the size written in header"
		}
	case TUserFileNotify:
		if msg.Header == nil || len(msg.Header) != 5 {
			errMsg = "header does not contain all required data"
		} else if msg.Header[0] == "" {
			errMsg = "sender name cannot be empty"
		} else if msg.Header[1] == "" {
			errMsg = "target name cannot be empty"
		} else if msg.Header[2] == "" {
			errMsg = "file name cannot be empty"
		} else if msg.Header[3] == "" {
			errMsg = "mime type cannot be empty"
		} else if bodySize, err := strconv.Atoi(msg.Header[4]); err != nil {
			errMsg = "header must contain the body size"
		} else if len(msg.Body) != bodySize {
			errMsg = "body does not match the size written in header"
		}
	case TError:
		if msg.Header == nil || len(msg.Header) != 1 {
			errMsg = "header does not contains all required data"
		} else if lines, err := strconv.Atoi(msg.Header[0]); err != nil {
			errMsg = "header must contain the body size"
		} else if bytes.Count(msg.Body, LineBreak) == lines {
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
