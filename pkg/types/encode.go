package types

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// EncodeMsg encodes given command and struct into bytes
// which are sent over the network. First byte indicates the command.
// Last byte is new line to indicate the end of the message.
func EncodeMsg(cmd CommandType, v interface{}) ([]byte, error) {
	msgBytes, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("could not marshal message to JSON: %v", err)
	}

	buf := new(bytes.Buffer)
	buf.WriteByte(byte(cmd))
	buf.Write(msgBytes)
	buf.WriteByte('\n')
	return buf.Bytes(), nil
}
