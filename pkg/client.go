package dslp

import (
	"log"
	"net"

	clientv1_2 "github.com/IceflowRE/go-dslp/pkg/client/v1_2"
	clientv2_0 "github.com/IceflowRE/go-dslp/pkg/client/v2_0"
)

func StartClient(address string, version string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Panicln(err)
	}
	log.Println("Connected to", address)
	switch version {
	case "1.2":
		go clientv1_2.HandleRequest(conn)
	case "2.0":
		go clientv2_0.HandleRequest(conn)
	default:
		return
	}
	defer conn.Close()
	defer log.Println("Closed connection to", address)

	switch version {
	case "1.2":
		clientv1_2.MainMenu(conn)
	case "2.0":
		clientv2_0.MainMenu(conn)
	default:
		return
	}
}
