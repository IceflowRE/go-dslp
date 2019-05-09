package client

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/IceflowRE/go-dslp/pkg/message"
	"github.com/IceflowRE/go-dslp/pkg/utils"
	"github.com/IceflowRE/go-dslp/pkg/v1_2"
)

var msgBuf = utils.NewCircularBuffer(50)

func HandleRequest(conn net.Conn) {
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
			msg, buf = v1_2.ScanMessage(buf)
			if msg != nil {
				err = msg.Valid()
				if err == nil {
					msgBuf.Add(msg)
				}
				if err != nil {
					conn.Write(v1_2.NewErrorMsg(err.Error()).ToBytes())
				}
			}

			// if valid message was found
			ok = msg != nil
		}

		if len(buf) > 16384 {
			conn.Write(v1_2.NewErrorMsg("Message exceeded 16384 bytes size. Disconnecting.").ToBytes())
			fmt.Println("Message exceeded 16384 bytes size. Disconnecting.")
			break
		}
	}
}
