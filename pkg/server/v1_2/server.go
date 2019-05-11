package serverv1_2

import (
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/IceflowRE/go-dslp/pkg/message"
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
			msg, buf = ScanMessage(buf)
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

func ScanMessage(data []byte) (message.IMessage, []byte) {
	res := rxMessage.FindSubmatchIndex(data)
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

// HandleMessage requires a valid message
func handleMessage(msg message.IMessage, conn net.Conn) error {
	switch msg.GetType() {
	case TRequestTime:
		message.SendMessage(conn, NewResponseTimeMsg())
	case TResponseTime:
		// do nothing
	case TGroupJoin:
		joinGroup(conn, *msg.GetContent())
	case TGroupLeave:
		return leaveGroup(conn, *msg.GetContent())
	case TGroupNotify:
		split := strings.SplitN(*msg.GetContent(), "\r\n", 2)
		sendToGroup(split[0], split[1])
	case TPeerNotify:
		split := strings.SplitN(*msg.GetContent(), "\r\n", 2)
		sendPeerNotify(net.ParseIP(split[0]), split[1])
	case TError:
		util.Println(conn, "Error message received", *msg.GetContent())
	default:
		util.Println(conn, "type invalid", msg.GetType())
	}
	return nil
}
