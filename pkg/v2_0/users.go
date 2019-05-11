package v2_0

import (
	"errors"
	"net"
	"sync"

	"github.com/IceflowRE/go-dslp/pkg/message"
	"github.com/IceflowRE/go-dslp/pkg/utils"
)

var nickConn = make(map[string]net.Conn)
var nickConnLock = sync.RWMutex{}
var connNick = make(map[net.Conn]string)
var connNickLock = sync.RWMutex{}

func registerUser(conn net.Conn, nick string) error {
	nickConnLock.Lock()
	connNickLock.Lock()
	defer nickConnLock.Unlock()
	defer connNickLock.Unlock()
	if _, ok := nickConn[nick]; ok {
		return errors.New("nick " + nick + " is already registered")
	}
	if _, ok := connNick[conn]; ok {
		return errors.New("you already registered a nick, unregister it first")
	}
	connNick[conn] = nick
	nickConn[nick] = conn
	utils.Println(conn, "NICK REGISTERED", nick)
	return nil
}

func unregisterUser(conn net.Conn, nick string) error {
	nickConnLock.Lock()
	connNickLock.Lock()
	defer nickConnLock.Unlock()
	defer connNickLock.Unlock()
	if tmp, ok := nickConn[nick]; !ok {
		return errors.New("nick " + nick + " is not registered")
	} else if tmp != conn {
		return errors.New("nick " + nick + " is not registered by you")
	}
	delete(connNick, conn)
	delete(nickConn, nick)
	utils.Println(conn, "NICK UNREGISTERED", nick)
	return nil
}

func freeUser(conn net.Conn) {
	nickConnLock.Lock()
	connNickLock.Lock()
	defer nickConnLock.Unlock()
	defer connNickLock.Unlock()
	if nick, ok := connNick[conn]; ok {
		delete(connNick, conn)
		delete(nickConn, nick)
	}
}

func nickExist(nick string) bool {
	nickConnLock.Lock()
	connNickLock.Lock()
	defer nickConnLock.Unlock()
	defer connNickLock.Unlock()
	_, ok := nickConn[nick]
	return ok
}

func nickValid(conn net.Conn, nick string) bool {
	nickConnLock.Lock()
	connNickLock.Lock()
	defer nickConnLock.Unlock()
	defer connNickLock.Unlock()
	regConn, ok := nickConn[nick]
	return ok && conn == regConn
}

// requires a valid message
func sendToUser(conn net.Conn, msg *Message) error {
	if !nickValid(conn, msg.Header[0]) {
		return errors.New("nick " + msg.Header[0] + " is not registered by you")
	} else if !nickExist(msg.Header[1]) {
		return errors.New("nick" + msg.Header[1] + " is not registered")
	}
	nickConnLock.Lock()
	defer nickConnLock.Unlock()
	trgtConn := nickConn[msg.Header[1]]
	message.SendMessage(trgtConn, msg)
	return nil
}
