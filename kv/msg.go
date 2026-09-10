package kv

import (
	"bytes"
	"encoding/binary"
	"net"
)

type Msg struct {
	Data []byte
	Conn net.Conn
}

func MsgLen(b []byte) int {
	return int(binary.LittleEndian.Uint32(b))
}

func ReadPayload(b []byte) []string {
	slices := bytes.Split(b, []byte{0})
	strs := make([]string, len(slices))
	for i, v := range slices {
		strs[i] = string(v)
	}
	return strs
}
