package messagev1_2

import "github.com/IceflowRE/go-dslp/pkg/message"

func ScanMessage(data []byte) (message.IMessage, []byte) {
	res := rxMessage.FindSubmatchIndex(data)
	if res != nil {
		msg := NewMessage()
		msg.Type = string(data[res[2]:res[3]])
		if res[4] != -1 {
			msg.Content = data[res[4]:res[5]]
		}
		return msg, data[res[1]:]
	}

	return nil, data
}
