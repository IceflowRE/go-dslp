package dslp

import (
	"log"
	"net"

	v12Client "github.com/IceflowRE/go-dslp/pkg/v1_2/client"
	v20Client "github.com/IceflowRE/go-dslp/pkg/v2_0/client"
)

func StartClient(address string, version string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Panicln(err)
	}
	log.Println("Connected to", address)
	switch version {
	case "1.2":
		go v12Client.HandleRequest(conn)
	case "2.0":
		go v20Client.HandleRequest(conn)
	default:
		return
	}
	defer conn.Close()
	defer log.Println("Closed connection to", address)

	switch version {
	case "1.2":
		v12Client.MainMenu(conn)
	case "2.0":
		v20Client.MainMenu(conn)
	default:
		return
	}
}
