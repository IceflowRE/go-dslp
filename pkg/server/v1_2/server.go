package serverv1_2

import (
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/IceflowRE/go-dslp/pkg/message"
	msgv1_2 "github.com/IceflowRE/go-dslp/pkg/message/v1_2"
	"github.com/IceflowRE/go-dslp/pkg/util"
)

func HandleRequest(conn net.Conn) {
	util.Println(conn, "accepted connection", "")
	addConn(conn)
	defer util.Println(conn, "closed connection", "")
	defer leaveAllGroups(conn)
	defer removeConn(conn)
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

		// do while message until valid message were received
		ok := true
		for ok {
			// TODO: maybe remove invalid buffer bytes
			var msg message.IMessage
			util.Println(conn, "BUFFER", buf)
			msg, buf = msgv1_2.ScanMessage(buf)
			if msg != nil {
				err = msg.Valid()
				if content := msg.GetContent(); content != nil {
					util.Println(conn, "RECEIVED ("+msg.GetType()+") valid: "+strconv.FormatBool(err == nil), *content)
				} else {
					util.Println(conn, "RECEIVED ("+msg.GetType()+") valid: "+strconv.FormatBool(err == nil), nil)
				}
				if err == nil {
					err = handleMessage(msg, conn)
				}
				if err != nil {
					message.SendMessage(conn, msgv1_2.NewErrorMsg(err.Error()))
					return
				}
			}

			// if valid message was found
			ok = msg != nil
		}

		if len(buf) > 16384 {
			message.SendMessage(conn, msgv1_2.NewErrorMsg("Message exceeded 16384 bytes size. Disconnecting."))
			return
		}
	}
}

// HandleMessage requires a valid message
func handleMessage(msg message.IMessage, conn net.Conn) error {
	switch msg.GetType() {
	case msgv1_2.TRequestTime:
		message.SendMessage(conn, msgv1_2.NewResponseTimeMsg())
	case msgv1_2.TResponseTime:
		// do nothing
	case msgv1_2.TGroupJoin:
		joinGroup(conn, *msg.GetContent())
	case msgv1_2.TGroupLeave:
		return leaveGroup(conn, *msg.GetContent())
	case msgv1_2.TGroupNotify:
		split := strings.SplitN(*msg.GetContent(), "\r\n", 2)
		sendToGroup(split[0], split[1])
	case msgv1_2.TPeerNotify:
		split := strings.SplitN(*msg.GetContent(), "\r\n", 2)
		sendPeerNotify(net.ParseIP(split[0]), split[1])
	case msgv1_2.TError:
		util.Println(conn, "Error message received", *msg.GetContent())
	default:
		util.Println(conn, "type invalid", msg.GetType())
	}
	return nil
}
