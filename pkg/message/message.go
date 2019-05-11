package message

import (
	"log"
	"net"

	"github.com/IceflowRE/go-dslp/pkg/util"
)

type IMessage interface {
	ToBytes() []byte
	GetType() string
	// return content of the message without the ending \r\n
	GetContent() *string
	GetRawContent() []byte
	Valid() error
}

func SendMessage(conn net.Conn, msg IMessage) {
	_, err := conn.Write(msg.ToBytes())
	if err != nil {
		log.Println(err)
	} else {
		util.Println(conn, "SENT ("+msg.GetType()+")", *msg.GetContent())
	}
}
