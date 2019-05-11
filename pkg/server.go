package dslp

import (
	"log"
	"net"
	"strconv"

	"github.com/IceflowRE/go-dslp/pkg/server/v1_2"
	"github.com/IceflowRE/go-dslp/pkg/server/v2_0"
)

func StartServer(port int, version string) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Panicln(err)
	}
	log.Println("Listening to connections on port", strconv.Itoa(port))
	defer listener.Close()

	var handlerFunc func(conn net.Conn)
	switch version {
	case "1.2":
		handlerFunc = serverv1_2.HandleRequest
	case "2.0":
		handlerFunc = serverv2_0.HandleRequest
	default:
		return
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Panicln(err)
		}
		go handlerFunc(conn)
	}
}
