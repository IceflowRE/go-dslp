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
	defer LeaveAllGroups(conn)
	defer conn.Close()

	buf := make([]byte, 0, 1024) // big buffer
	tmp := make([]byte, 256)     // small buffer
	for {
		n, err := conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				log.Println("read error:", err)
			}
			break
		}
		buf = append(buf, tmp[:n]...)

		// work until all valid message are proceeded
		ok := true
		for ok {
			// TODO: maybe remove invalid buffer bytes
			var msg message.IMessage
			utils.Println(conn, "BUFFER", buf)
			msg, buf = ScanMessage(buf)
			if msg != nil {
				err = msg.Valid()
				if msg.GetContent() != nil {
					utils.Println(conn, "RECEIVED ("+msg.GetType()+") valid: "+strconv.FormatBool(err == nil), *msg.GetContent())
				} else {
					utils.Println(conn, "RECEIVED ("+msg.GetType()+") valid: "+strconv.FormatBool(err == nil), nil)
				}
				if err == nil {
					err = HandleMessage(msg, conn)
				}
				if err != nil {
					conn.Write(NewErrorMsg(err.Error()).ToBytes())
				}
			}

			// if valid message was found
			ok = msg != nil
		}

		if len(buf) > 16384 {
			conn.Write(NewErrorMsg("Message exceeded 16384 bytes size. Disconnecting.").ToBytes())
			break
		}
	}
}

func ScanMessage(data []byte) (message.IMessage, []byte) {
	res := RxHeader.FindSubmatchIndex(data)
	if res == nil {
		return nil, data
	}
	// message end position
	msgEnd := res[1]

	msg := NewMessage()
	msg.Type = string(data[res[2]:res[3]])
	// if a header exists
	if res[4] != -1 {
		msg.Header = strings.Split(string(data[res[4]:res[5]]), "\r\n")

		// if message with body size information
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
		default:
			return msg, data[msgEnd:]
		}
		// if header is too small, return invalid message
		if len(msg.Header) < sizePos {
			return msg, data[msgEnd:]
		}
		bodySize, err := strconv.Atoi(msg.Header[sizePos])
		// header is malformed, return invalid message
		if err != nil {
			return msg, data[msgEnd:]
		}

		// get body
		switch msg.Type {
		case TGroupNotify, TUserTextNotify, TError:
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
	}
	return msg, data[msgEnd:]
}
