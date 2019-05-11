package clientv2_0

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	msgv2_0 "github.com/IceflowRE/go-dslp/pkg/message/v2_0"
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
		for idx, msgType := range msgv2_0.Types {
			fmt.Println(strconv.Itoa(idx) + ": " + msgType)
		}
		fmt.Println(strconv.Itoa(len(msgv2_0.Types)) + ": Back")

		input.Scan()
		selec, err := strconv.Atoi(input.Text())
		if err == nil && selec < len(msgv2_0.Types) {
			writeMessage(conn, msgv2_0.Types[selec])
			validInput = true
		} else if err == nil && selec == len(msgv2_0.Types) {
			break
		}
	}
}

func writeMessage(conn net.Conn, msgType string) {
	var err error
	switch msgType {
	case msgv2_0.TRequestTime:
		_, err = conn.Write(msgv2_0.NewRequestTime().ToBytes())
	case msgv2_0.TResponseTime:
		_, err = conn.Write(msgv2_0.NewResponseTimeMsg().ToBytes())
	case msgv2_0.TGroupJoin:
		fmt.Println("Group to join:")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()

		_, err = conn.Write(msgv2_0.NewGroupJoin(input.Text()).ToBytes())
	case msgv2_0.TGroupLeave:
		fmt.Println("Group to leave:")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()

		_, err = conn.Write(msgv2_0.NewGroupLeave(input.Text()).ToBytes())
	case msgv2_0.TGroupNotify:
		fmt.Println("Group to notify:")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()

		group := input.Text()

		fmt.Println("Message (empty line to end):")
		content := getContent()

		_, err = conn.Write(msgv2_0.NewGroupNotify(group, content).ToBytes())
	case msgv2_0.TUserJoin:
		fmt.Println("username to register:")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()

		_, err = conn.Write(msgv2_0.NewUserJoin(input.Text()).ToBytes())
	case msgv2_0.TUserLeave:
		fmt.Println("username to unregister:")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()

		_, err = conn.Write(msgv2_0.NewUserLeave(input.Text()).ToBytes())
	case msgv2_0.TUserTextNotify:
		input := bufio.NewScanner(os.Stdin)
		fmt.Println("username to send:")
		input.Scan()
		sender := input.Text()

		fmt.Println("user to notify:")
		input.Scan()
		target := input.Text()

		fmt.Println("Message (empty line to end):")
		content := getContent()

		_, err = conn.Write(msgv2_0.NewUserTextNotify(sender, target, content).ToBytes())
	case msgv2_0.TUserFileNotify:
		input := bufio.NewScanner(os.Stdin)
		fmt.Println("username to send:")
		input.Scan()
		sender := input.Text()

		fmt.Println("user to notify:")
		input.Scan()
		target := input.Text()

		fmt.Println("filename:")
		input.Scan()
		filename := input.Text()

		var file []byte
		file, err = ioutil.ReadFile(filename) // b has type []byte
		if err != nil {
			err = errors.New("cannot use filename: " + err.Error())
			break
		}

		var contentType string
		if contentType = http.DetectContentType(file[:512]); contentType == "application/octet-stream" {
			fmt.Println("could not detect mime type automatically, please specify yourself:")
			input.Scan()
			contentType = input.Text()
		}

		_, err = conn.Write(msgv2_0.NewUserFileNotify(sender, target, filepath.Base(filename), contentType, file).ToBytes())
	case msgv2_0.TError:
		fmt.Println("Message (empty line to end):")
		content := getContent()

		_, err = conn.Write(msgv2_0.NewErrorMsg(content[0]).ToBytes())
	default:
		err = errors.New("Cannot handle the chosen message type.")
	}
	if err != nil {
		fmt.Println("Error: " + err.Error())
	} else {
		fmt.Println("Message send successfully.")
	}
}

func getContent() []string {
	input := bufio.NewScanner(os.Stdin)

	msg := make([]string, 0)
	for input.Scan() {
		if input.Text() == "" {
			break
		} else {
			msg = append(msg, input.Text())
		}
	}
	return msg
}

func showMessages() {
	fmt.Println(strings.Repeat("=", 7), "New Messages", strings.Repeat("=", 7))
	for msgBuf.Size() > 0 {
		msg := msgBuf.Remove().(*msgv2_0.Message)

		if msg.Type == msgv2_0.TUserFileNotify {
			err := ioutil.WriteFile(msg.Header[3], msg.Body, 0644)
			if err == nil {
				fmt.Println(msg.GetType(), " || ", "type (", msg.Header[2], "), size (", msg.Header[4], ") saved as", msg.Header[3])
			} else {
				fmt.Println(msg.GetType(), " || ", "type (", msg.Header[2], "), size (", msg.Header[4], ") name(", msg.Header[3], ") could not be saved: ", err.Error())
			}
		} else {
			content := msg.GetContent()
			if content == nil {
				fmt.Println(msg.GetType(), " || ", content)
			} else {
				fmt.Println(msg.GetType(), " || ", *content)
			}
		}
	}
	fmt.Println(strings.Repeat("=", 28))
}
