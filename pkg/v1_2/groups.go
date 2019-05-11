package v1_2

import (
	"errors"
	"net"
	"sync"

	"github.com/IceflowRE/go-dslp/pkg/message"
	"github.com/IceflowRE/go-dslp/pkg/utils"
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
	utils.Println(conn, "GROUP JOIN", group)
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
			utils.Println(conn, "GROUP LEAVE", group)
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
			utils.Println(conn, "GROUP LEAVE", group)
			return nil
		}
	}
	return errors.New("you are not a member of this group")
}

func sendToGroup(group string, content string) {
	groupsLock.RLock()
	defer groupsLock.RUnlock()
	msg := NewGroupNotify(group, content)
	if value, ok := groups[group]; ok {
		for member := range value {
			message.SendMessage(member, msg)
		}
	}
}
