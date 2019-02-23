package types

type CommandType byte

const (
	CmdLogin        CommandType = 0x1
	CmdPing         CommandType = 0x2
	CmdStatusChange CommandType = 0x3
)

type Msg struct {
	Cmd  CommandType
	Data []byte
}

type LoginRequest struct {
	UserID  int   `json:"user_id"`
	Friends []int `json:"friends"`
}

type PingRequest struct {
	UserID int `json:"user_id"`
}

type StatusChangeReply struct {
	UserID int  `json:"user_id"`
	Online bool `json:"online"`
}
