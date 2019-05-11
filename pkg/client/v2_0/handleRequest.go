package clientv2_0

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/IceflowRE/go-dslp/pkg/message"
	msgv2_0 "github.com/IceflowRE/go-dslp/pkg/message/v2_0"
	"github.com/IceflowRE/go-dslp/pkg/util"
)

var msgBuf = util.NewCircularBuffer(50)

func HandleRequest(conn net.Conn) {
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
			var msg *msgv2_0.Message
			msg, buf = msgv2_0.ScanMessage(buf)
			if msg != nil {
				err = msg.Valid()
				if err == nil {
					if msg.Type == msgv2_0.TError {
						fmt.Println("Error message was received and the connection may closed.")
					}
					msgBuf.Add(msg)
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
			message.SendMessage(conn, msgv2_0.NewErrorMsg("Message exceeded 16384 bytes size. Disconnecting"))
			fmt.Println("Message exceeded 16384 bytes size. Disconnecting.")
			return
		}
	}
}
