package serverv2_0

import (
	"errors"
	"net"
	"sync"

	"github.com/IceflowRE/go-dslp/pkg/message"
	"github.com/IceflowRE/go-dslp/pkg/util"
)

var groups = make(map[string]map[net.Conn]struct{})
var groupsLock = sync.RWMutex{}

func joinGroup(conn net.Conn, group string) {
	groupsLock.Lock()
	defer groupsLock.Unlock()
	if _, ok := groups[group]; !ok {
		groups[group] = make(map[net.Conn]struct{})
	}
	groups[group][conn] = struct{}{}
	util.Println(conn, "GROUP JOIN", group)
}

func leaveAllGroups(conn net.Conn) {
	groupsLock.Lock()
	defer groupsLock.Unlock()
	for group, value := range groups {
		if _, ok := value[conn]; ok {
			delete(value, conn)
			if len(value) == 0 {
				delete(groups, group)
			}
			util.Println(conn, "GROUP LEAVE", group)
		}
	}
}

func leaveGroup(conn net.Conn, group string) error {
	groupsLock.Lock()
	defer groupsLock.Unlock()
	if value, ok := groups[group]; ok {
		if _, ok := value[conn]; ok {
			delete(value, conn)
			if len(value) == 0 {
				delete(groups, group)
			}
			util.Println(conn, "GROUP LEAVE", group)
			return nil
		}
	}
	return errors.New("you are not a member of group " + group)
}

// requires a valid message
func sendToGroup(conn net.Conn, msg *Message) error {
	groupsLock.RLock()
	defer groupsLock.RUnlock()
	if value, ok := groups[msg.Header[0]]; ok {
		for member := range value {
			// do not send to ourself again
			//if member != conn {
			message.SendMessage(member, msg)
			//}
		}
	} else {
		return errors.New("you are not a member of group " + msg.Header[0])
	}
	return nil
}
