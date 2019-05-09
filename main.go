package main

import (
	"flag"
	"github.com/IceflowRE/go-dslp/pkg"
	"strconv"
)

func main() {
	server := flag.String("server", "", "<port>  | StartServer server.")
	client := flag.String("client", "", "<address>  | StartServer client.")
	version := flag.String("version", "", "<version>  | Use this protocol version. Available: 1.2 | 2.0")
	flag.Parse()
	if server != nil && version != nil {
		port, err := strconv.Atoi(*server)
		if err == nil {
			dslp.StartServer(port, *version)
		} else {
			flag.Usage()
		}
	} else if client != nil && version != nil {
		dslp.StartClient(*client, *version)
	} else {
		flag.Usage()
	}
}
