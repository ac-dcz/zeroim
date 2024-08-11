package protocol

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
)

var (
	ErrHeaderInVaild = func(size int) error {
		return fmt.Errorf("header size is %d ,no equal %d", size, headersize)
	}
	ErrPackSizeLength = func(n int) error {
		return fmt.Errorf("read package length: n=%d len=%d", n, packsizeLen)
	}
	ErrWriteLength = func(n, buf int) error {
		return fmt.Errorf("write conn: n=%d len(buf)=%d", n, buf)
	}
	ErrReadLength = func(n, buf int) error {
		return fmt.Errorf("read conn: n=%d len(buf)=%d", n, buf)
	}
)

type IMMessage struct {
	H    Header  `zeroim:"header"`
	Body Message `zeroim:"body"`
}

const (
	versionOffset    int = 0
	statusCodeOffset int = 1
	msgTypeOffset    int = 3
	seqOffset        int = 5
)

const (
	versionLen    int = 1
	statusCodeLen int = 2
	msgTypeLen    int = 2
	seqLen        int = 4
	packsizeLen   int = 4
	headersizeLen int = 2
)

const headersize = versionLen + statusCodeLen + msgTypeLen + seqLen

type Header struct {
	Version    uint8  `zeroim:"version"`     //服务端版本号
	StatusCode uint16 `zeroim:"status_code"` //状态码
	MsgType    uint16 `zeroim:"msg_type"`    //消息类型
	Seq        uint32 `zeroim:"req_seq"`     //消息序号
}

func (h *Header) Encode() ([]byte, error) {
	data := make([]byte, headersize)
	data[versionOffset] = h.Version
	binary.BigEndian.PutUint16(data[statusCodeOffset:], h.StatusCode)
	binary.BigEndian.PutUint16(data[msgTypeOffset:], h.MsgType)
	binary.BigEndian.PutUint32(data[seqOffset:], h.Seq)
	return data, nil
}

func (h *Header) Decode(data []byte) error {
	if len(data) != headersize {
		return ErrHeaderInVaild(len(data))
	}
	h.Version = data[versionOffset]
	h.StatusCode = binary.BigEndian.Uint16(data[statusCodeOffset : statusCodeOffset+statusCodeLen])
	h.MsgType = binary.BigEndian.Uint16(data[msgTypeOffset : msgTypeOffset+msgTypeLen])
	h.Seq = binary.BigEndian.Uint32(data[seqOffset:])
	return nil
}

type Message interface {
	Encode() ([]byte, error)
	Decode(data []byte) error
	MsgType() uint16
}

type PrivateChatMessage struct {
	From        uint64 `zeroim:"from"`
	To          uint64 `zeroim:"to"`
	Ts          uint64 `zeroim:"ts"`
	ContentType uint64 `zeroim:"content_type"`
	Data        []byte `zeroim:"data"`
}

func (m *PrivateChatMessage) Encode() ([]byte, error) {
	return json.Marshal(m)
}

func (m *PrivateChatMessage) Decode(data []byte) error {
	return json.Unmarshal(data, m)
}

func (m *PrivateChatMessage) MsgType() uint16 {
	return PrivateChat
}

type RoomChatMessage struct {
	From        uint64 `zeroim:"from"`
	RoomID      uint64 `zeroim:"room_id"`
	Ts          uint64 `zeroim:"ts"`
	ContentType uint64 `zeroim:"content_type"`
	Data        []byte `zeroim:"data"`
}

func (m *RoomChatMessage) Encode() ([]byte, error) {
	return json.Marshal(m)
}

func (m *RoomChatMessage) Decode(data []byte) error {
	return json.Unmarshal(data, m)
}

func (m *RoomChatMessage) MsgType() uint16 {
	return RoomChat
}

const (
	PrivateChat uint16 = iota
	RoomChat
)
