package clientv1_2

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	msgv1_2 "github.com/IceflowRE/go-dslp/pkg/message/v1_2"
)

func MainMenu(conn net.Conn) {
	input := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(
			"DSLP Client:\n",
			"0: Show new Messages (50 max.)\n",
			"1: Write Message\n",
			"2: Exit\n",
		)
		input.Scan()
		switch input.Text() {
		case "0":
			showMessages()
		case "1":
			writeMenu(conn)
		case "2":
			return
		}
	}
}

func writeMenu(conn net.Conn) {
	validInput := false
	input := bufio.NewScanner(os.Stdin)
	for !validInput {
		fmt.Println("Write Message:")
		for idx, msgType := range msgv1_2.Types {
			fmt.Println(strconv.Itoa(idx) + ": " + msgType)
		}
		fmt.Println(strconv.Itoa(len(msgv1_2.Types)) + ": Back")

		input.Scan()
		selec, err := strconv.Atoi(input.Text())
		if err == nil && selec < len(msgv1_2.Types) {
			writeMessage(conn, msgv1_2.Types[selec])
			validInput = true
		} else if err == nil && selec == len(msgv1_2.Types) {
			break
		}
	}
}

func writeMessage(conn net.Conn, msgType string) {
	var err error
	switch msgType {
	case msgv1_2.TRequestTime:
		_, err = conn.Write(msgv1_2.NewRequestTime().ToBytes())
	case msgv1_2.TResponseTime:
		_, err = conn.Write(msgv1_2.NewResponseTimeMsg().ToBytes())
	case msgv1_2.TGroupJoin:
		fmt.Println("Group to join:")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()

		_, err = conn.Write(msgv1_2.NewGroupJoin(input.Text()).ToBytes())
	case msgv1_2.TGroupLeave:
		fmt.Println("Group to leave:")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()

		_, err = conn.Write(msgv1_2.NewGroupLeave(input.Text()).ToBytes())
	case msgv1_2.TGroupNotify:
		fmt.Println("Group to notify:")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()

		group := input.Text()

		fmt.Println("Message (empty line to end):")
		content := getContent()

		_, err = conn.Write(msgv1_2.NewGroupNotify(group, content).ToBytes())
	case msgv1_2.TPeerNotify:
		fmt.Println("IP to notify:")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()

		ip := net.ParseIP(input.Text())
		if ip == nil {
			err = errors.New("IP was not in a valid format.")
			break
		}

		fmt.Println("Message (empty line to end):")
		content := getContent()

		_, err = conn.Write(msgv1_2.NewPeerNotfiy(ip, content).ToBytes())
	case msgv1_2.TError:
		fmt.Println("Message (empty line to end):")
		content := getContent()

		_, err = conn.Write(msgv1_2.NewErrorMsg(content).ToBytes())
	default:
		err = errors.New("Cannot handle the chosen message type.")
	}
	if err != nil {
		fmt.Println("Error: " + err.Error())
	} else {
		fmt.Println("Message send successfully.")
	}
}

func getContent() string {
	input := bufio.NewScanner(os.Stdin)

	msg := make([]string, 0)
	for input.Scan() {
		if input.Text() == "" {
			break
		} else {
			msg = append(msg, input.Text())
		}
	}
	return strings.Join(msg, "\r\n")
}

func showMessages() {
	fmt.Println(strings.Repeat("=", 7), "New Messages", strings.Repeat("=", 7))
	for msgBuf.Size() > 0 {
		msg := msgBuf.Remove().(*msgv1_2.Message)

		content := msg.GetContent()
		if content == nil {
			fmt.Println(msg.GetType(), " || ", content)
		} else {
			fmt.Println(msg.GetType(), " || ", *content)
		}
	}
	fmt.Println(strings.Repeat("=", 28))
}
