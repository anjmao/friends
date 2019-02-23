package types

import (
	"fmt"
	"reflect"
	"testing"
)

func TestEncodeMessages(t *testing.T) {
	tests := []struct {
		cmd            CommandType
		v              interface{}
		expectedOutput string
	}{
		{
			cmd: CmdLogin,
			v: LoginRequest{
				UserID:  1,
				Friends: []int{2},
			},
			expectedOutput: "017b22757365725f6964223a312c22667269656e6473223a5b325d7d0a",
		},
		{
			cmd: CmdPing,
			v: PingRequest{
				UserID: 1,
			},
			expectedOutput: "027b22757365725f6964223a317d0a",
		},
		{
			cmd: CmdStatusChange,
			v: StatusChangeReply{
				UserID: 1,
				Online: true,
			},
			expectedOutput: "037b22757365725f6964223a312c226f6e6c696e65223a747275657d0a",
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("encode %b", test.cmd), func(tt *testing.T) {
			b, err := EncodeMsg(test.cmd, test.v)
			if err != nil {
				tt.Fatalf("could not encode message: %v", err)
			}
			bytesHex := fmt.Sprintf("%x", b)
			if !reflect.DeepEqual(bytesHex, test.expectedOutput) {
				tt.Fatalf("expected output %s, got %s", test.expectedOutput, bytesHex)
			}
		})
	}
}
