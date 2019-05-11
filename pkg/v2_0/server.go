package v2_0

import (
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/IceflowRE/go-dslp/pkg/message"
	"github.com/IceflowRE/go-dslp/pkg/utils"
)

func HandleRequest(conn net.Conn) {
	utils.Println(conn, "accepted connection", "")
	defer utils.Println(conn, "closed connection", "")
	defer leaveAllGroups(conn)
	defer freeUser(conn)
	defer conn.Close()

	buf := make([]byte, 0, 1024) // big buffer
	tmp := make([]byte, 256)     // small buffer
	for {
		n, err := conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				log.Println("read error:", err)
			}
			return
		}
		buf = append(buf, tmp[:n]...)

		// work until all valid message are proceeded
		ok := true
		for ok {
			utils.Println(conn, "BUFFER", buf)
			var msg *Message
			msg, buf = ScanMessage(buf)
			if msg != nil {
				err = msg.Valid()
				if msg.GetContent() != nil {
					utils.Println(conn, "RECEIVED ("+msg.GetType()+") valid: "+strconv.FormatBool(err == nil), *msg.GetContent())
				} else {
					utils.Println(conn, "RECEIVED ("+msg.GetType()+") valid: "+strconv.FormatBool(err == nil), nil)
				}
				if err == nil {
					err = handleMessage(msg, conn)
				}
				if err != nil {
					message.SendMessage(conn, NewErrorMsg(err.Error()))
					return
				}
			}

			// if valid message was found
			ok = msg != nil
		}

		if len(buf) > 16384 {
			message.SendMessage(conn, NewErrorMsg("Message exceeded 16384 bytes size. Disconnecting."))
			return
		}
	}
}

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
		bodyEndPos := utils.IndexN(data[res[1]:], LineBreak, bodySize)
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

// handleMessage requires a valid message
func handleMessage(msg *Message, conn net.Conn) error {
	switch msg.GetType() {
	case TRequestTime:
		message.SendMessage(conn, NewResponseTimeMsg())
	case TResponseTime:
		// do nothing
	case TGroupJoin:
		joinGroup(conn, msg.Header[0])
	case TGroupLeave:
		return leaveGroup(conn, *msg.GetContent())
	case TGroupNotify:
		return sendToGroup(conn, msg)
	case TUserJoin:
		return registerUser(conn, msg.Header[0])
	case TUserLeave:
		return unregisterUser(conn, msg.Header[0])
	case TUserTextNotify, TUserFileNotify:
		return sendToUser(conn, msg)
	case TError:
		utils.Println(conn, "Error message received", *msg.GetContent())
	default:
		utils.Println(conn, "type invalid", msg.GetType())
	}
	return nil
}
