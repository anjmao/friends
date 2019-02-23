package types

// DecodeMsg decodes bytes into Msg struct.
func DecodeMsg(b []byte) *Msg {
	if len(b) == 0 {
		return &Msg{}
	}
	cmd := CommandType(b[0])
	data := b[1:]
	return &Msg{Cmd: cmd, Data: data}
}
