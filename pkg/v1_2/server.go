package v1_2

import (
	"io"
	"log"
	"net"
	"strconv"

	"github.com/IceflowRE/go-dslp/pkg/message"
	"github.com/IceflowRE/go-dslp/pkg/utils"
)

func HandleRequest(conn net.Conn) {
	utils.Println(conn, "accepted connection", "")
	AddConn(conn)
	defer utils.Println(conn, "closed connection", "")
	defer LeaveAllGroups(conn)
	defer RemoveConn(conn)
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

		// do while message until valid message were received
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
	res := RxMessage.FindSubmatchIndex(data)
	if res != nil {
		msg := NewMessage()
		msg.Type = string(data[res[2]:res[3]])
		if res[4] != -1 {
			msg.Content = data[res[4]:res[5]]
		}

		return msg, data[res[1]:]
	}

	return nil, data
}
