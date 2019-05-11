package serverv1_2

import (
	"net"
	"sync"

	"github.com/IceflowRE/go-dslp/pkg/message"
	msgv1_2 "github.com/IceflowRE/go-dslp/pkg/message/v1_2"
)

var connections = make(map[string]net.Conn)
var connectionsLock = sync.RWMutex{}

func addConn(conn net.Conn) {
	connectionsLock.Lock()
	defer connectionsLock.Unlock()
	connections[conn.RemoteAddr().(*net.TCPAddr).IP.String()] = conn
}

func removeConn(conn net.Conn) {
	connectionsLock.Lock()
	defer connectionsLock.Unlock()
	delete(connections, conn.RemoteAddr().(*net.TCPAddr).IP.String())
}

func sendPeerNotify(ip net.IP, content string) {
	connectionsLock.RLock()
	defer connectionsLock.RUnlock()
	ipCmp := ip.String()
	msg := msgv1_2.NewPeerNotfiy(ip, content)
	for ipStr, conn := range connections {
		if ipStr == ipCmp {
			message.SendMessage(conn, msg)
			break
		}
	}
}
