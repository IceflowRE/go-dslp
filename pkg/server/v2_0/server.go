package serverv2_0

import (
	"io"
	"log"
	"net"
	"strconv"

	"github.com/IceflowRE/go-dslp/pkg/message"
	msgv2_0 "github.com/IceflowRE/go-dslp/pkg/message/v2_0"
	"github.com/IceflowRE/go-dslp/pkg/util"
)

func HandleRequest(conn net.Conn) {
	util.Println(conn, "accepted connection", "")
	defer util.Println(conn, "closed connection", "")
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
			util.Println(conn, "BUFFER", buf)
			var msg *msgv2_0.Message
			msg, buf = msgv2_0.ScanMessage(buf)
			if msg != nil {
				err = msg.Valid()
				if msg.GetContent() != nil {
					util.Println(conn, "RECEIVED ("+msg.GetType()+") valid: "+strconv.FormatBool(err == nil), *msg.GetContent())
				} else {
					util.Println(conn, "RECEIVED ("+msg.GetType()+") valid: "+strconv.FormatBool(err == nil), nil)
				}
				if err == nil {
					err = handleMessage(msg, conn)
				}
				if err != nil {
					message.SendMessage(conn, msgv2_0.NewErrorMsg(err.Error()))
					return
				}
			}

			// if valid message was found
			ok = msg != nil
		}

		if len(buf) > 16384 {
			message.SendMessage(conn, msgv2_0.NewErrorMsg("Message exceeded 16384 bytes size. Disconnecting."))
			return
		}
	}
}

// handleMessage requires a valid message
func handleMessage(msg *msgv2_0.Message, conn net.Conn) error {
	switch msg.GetType() {
	case msgv2_0.TRequestTime:
		message.SendMessage(conn, msgv2_0.NewResponseTimeMsg())
	case msgv2_0.TResponseTime:
		// do nothing
	case msgv2_0.TGroupJoin:
		joinGroup(conn, msg.Header[0])
	case msgv2_0.TGroupLeave:
		return leaveGroup(conn, *msg.GetContent())
	case msgv2_0.TGroupNotify:
		return sendToGroup(conn, msg)
	case msgv2_0.TUserJoin:
		return registerUser(conn, msg.Header[0])
	case msgv2_0.TUserLeave:
		return unregisterUser(conn, msg.Header[0])
	case msgv2_0.TUserTextNotify, msgv2_0.TUserFileNotify:
		return sendToUser(conn, msg)
	case msgv2_0.TError:
		util.Println(conn, "Error message received", *msg.GetContent())
	default:
		util.Println(conn, "type invalid", msg.GetType())
	}
	return nil
}
