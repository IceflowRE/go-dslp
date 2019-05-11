package main

import (
	"flag"
	"strconv"

	"github.com/IceflowRE/go-dslp/pkg"
)

func main() {
	server := flag.String("server", "", "<port>  | StartServer server.")
	client := flag.String("client", "", "<address>  | StartServer client.")
	version := flag.String("version", "", "<version>  | Use this protocol version. Available: 1.2 | 2.0")
	flag.Parse()
	if *server != "" && *version != "" {
		port, err := strconv.Atoi(*server)
		if err == nil {
			dslp.StartServer(port, *version)
		} else {
			flag.Usage()
		}
	} else if *client != "" && *version != "" {
		dslp.StartClient(*client, *version)
	} else {
		flag.Usage()
	}
}
