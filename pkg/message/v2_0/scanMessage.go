package messagev2_0

import (
	"strconv"
	"strings"

	"github.com/IceflowRE/go-dslp/pkg/util"
)

func ScanMessage(data []byte) (*Message, []byte) {
	res := RxHeader.FindSubmatchIndex(data)
	if res == nil {
		return nil, data
	}
	// message end position
	msgEnd := res[1]

	msg := NewMessage()
	msg.Type = string(data[res[2]:res[3]])

	bodySize := 0
	// if a header exists
	if res[4] != -1 {
		msg.Header = strings.Split(string(data[res[4]:res[5]]), "\r\n")

		// if message with body size information, get the body size
		switch msg.Type {
		case TGroupNotify, TUserTextNotify, TUserFileNotify, TError:
			// where in the header the size is written
			sizePos := -1
			switch msg.Type {
			case TGroupNotify:
				sizePos = 1
			case TUserTextNotify:
				sizePos = 2
			case TUserFileNotify:
				sizePos = 4
			case TError:
				sizePos = 0
			}
			// if header is too small, return invalid message
			if len(msg.Header) < sizePos {
				return msg, data[msgEnd:]
			}
			var err error
			bodySize, err = strconv.Atoi(msg.Header[sizePos])
			// header is malformed, return invalid message
			if err != nil {
				return msg, data[msgEnd:]
			}
		}
	}
	if msg.Type == TResponseTime {
		bodySize = 1
	}

	// get body
	switch msg.Type {
	case TResponseTime, TGroupNotify, TUserTextNotify, TError:
		bodyEndPos := util.IndexN(data[res[1]:], LineBreak, bodySize)
		// if not the whole body is available wait for the missing data
		if bodyEndPos == -1 {
			return nil, data
		}

		// add body length and linebreak length to message end index
		msgEnd += bodyEndPos + len(LineBreak)
		msg.Body = data[res[1]:msgEnd]
	case TUserFileNotify:
		// if not the whole body is available wait for the missing data
		if len(data) < msgEnd+bodySize {
			return nil, data
		}
		msgEnd += bodySize
		msg.Body = data[res[1]:msgEnd]
	}
	return msg, data[msgEnd:]
}
