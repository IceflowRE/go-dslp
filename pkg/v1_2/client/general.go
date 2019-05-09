package client

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/IceflowRE/go-dslp/pkg/v1_2"
)

func Start(address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Panicln(err)
	}
	log.Println("Connected to", address)
	go HandleRequest(conn)
	defer conn.Close()
	defer log.Println("Closed connection to", address)
	MainMenu(conn)
}

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
		for idx, msgType := range v1_2.Types {
			fmt.Println(strconv.Itoa(idx) + ": " + msgType)
		}
		fmt.Println(strconv.Itoa(len(v1_2.Types)) + ": Back")

		input.Scan()
		selec, err := strconv.Atoi(input.Text())
		if err == nil && selec < len(v1_2.Types) {
			writeMessage(conn, v1_2.Types[selec])
			validInput = true
		} else if err == nil && selec == len(v1_2.Types) {
			break
		}
	}
}

func writeMessage(conn net.Conn, msgType string) {
	var err error
	switch msgType {
	case v1_2.TRequestTime:
		_, err = conn.Write(v1_2.NewRequestTime().ToBytes())
	case v1_2.TResponseTime:
		_, err = conn.Write(v1_2.NewResponseTimeMsg().ToBytes())
	case v1_2.TGroupJoin:
		fmt.Println("Group to join:")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()

		_, err = conn.Write(v1_2.NewGroupJoin(input.Text()).ToBytes())
	case v1_2.TGroupLeave:
		fmt.Println("Group to leave:")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()

		_, err = conn.Write(v1_2.NewGroupLeave(input.Text()).ToBytes())
	case v1_2.TGroupNotify:
		fmt.Println("Group to notify:")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()

		group := input.Text()

		fmt.Println("Message (empty line to end):")
		content := getContent()

		_, err = conn.Write(v1_2.NewGroupNotify(group, content).ToBytes())
	case v1_2.TPeerNotify:
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

		_, err = conn.Write(v1_2.NewPeerNotfiy(ip, content).ToBytes())
	case v1_2.TError:
		fmt.Println("Message (empty line to end):")
		content := getContent()

		_, err = conn.Write(v1_2.NewErrorMsg(content).ToBytes())
	default:
		err = errors.New("Cannot handle the choosen message type.")
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
		msg := msgBuf.Remove().(*v1_2.Message)

		content := msg.GetContent()
		if content == nil {
			fmt.Println(msg.GetType(), " || ", content)
		} else {
			fmt.Println(msg.GetType(), " || ", *content)
		}
	}
	fmt.Println(strings.Repeat("=", 28))
}
