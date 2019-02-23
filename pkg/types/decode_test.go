package types

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"
)

func TestDecodeMessages(t *testing.T) {
	tests := []struct {
		cmd            CommandType
		input          string
		expectedOutput string
	}{
		{
			cmd:            CmdLogin,
			input:          "017b22757365725f6964223a312c22667269656e6473223a5b325d7d0a",
			expectedOutput: "{\"Cmd\":1,\"Data\":\"eyJ1c2VyX2lkIjoxLCJmcmllbmRzIjpbMl19Cg==\"}",
		},
		{
			cmd:            CmdPing,
			input:          "027b22757365725f6964223a317d0a",
			expectedOutput: "{\"Cmd\":2,\"Data\":\"eyJ1c2VyX2lkIjoxfQo=\"}",
		},
		{
			cmd:            CmdStatusChange,
			input:          "037b22757365725f6964223a312c226f6e6c696e65223a747275657d0a",
			expectedOutput: "{\"Cmd\":3,\"Data\":\"eyJ1c2VyX2lkIjoxLCJvbmxpbmUiOnRydWV9Cg==\"}",
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("decode %b", test.cmd), func(tt *testing.T) {
			inputBytes, _ := hex.DecodeString(test.input)
			msg := DecodeMsg(inputBytes)
			jsonStr, _ := json.Marshal(msg)

			if string(jsonStr) != test.expectedOutput {
				tt.Fatalf("expected output %s, got %s", test.expectedOutput, jsonStr)
			}
		})
	}
}
