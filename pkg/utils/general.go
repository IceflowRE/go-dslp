package utils

import (
	"bytes"
	"log"
	"net"
	"strings"
)

func Println(conn net.Conn, tag string, msg interface{}) {
	log.Println(conn.RemoteAddr().String(), "||", strings.ToUpper(tag), "|", msg)
}

// Find the nth index of the separator.
// Returns the beginning of the nth separator.
// Returns -1 if the nth separator was not found.
func IndexN(s []byte, sep []byte, n int) int {
	i := 0
	c := 0
	if n <= 0 {
		return -1
	}
	for cur := 0; cur < n; cur++ {
		c = bytes.Index(s[c:], sep)
		if c == -1 {
			return -1
		} else {
			i += c
		}
	}
	return i
}
