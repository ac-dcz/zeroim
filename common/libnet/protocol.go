package libnet

import (
	"encoding/binary"
	"fmt"
	"net"
)

var (
	ErrPackSizeLength = func(n int) error {
		return fmt.Errorf("read package length: n=%d len=%d", n, PackSize)
	}
	ErrWriteLength = func(n, buf int) error {
		return fmt.Errorf("write conn: n=%d len(buf)=%d", n, buf)
	}
	ErrReadLength = func(n, buf int) error {
		return fmt.Errorf("read conn: n=%d len(buf)=%d", n, buf)
	}
)

const HeaderLength = 10
const (
	PackSize   = 4
	HeaderSize = 2
)
const (
	VersionOffset   = 0
	StatusOffset    = 1
	ServiceIdOffset = 2
	CmdOffset       = 4
	SeqOffset       = 6
)

type Header struct {
	Version   uint8
	Status    uint8
	ServiceId uint16
	Cmd       uint16
	Seq       uint32
}

func (h *Header) encode() []byte {
	buf := make([]byte, HeaderLength)
	buf[VersionOffset] = h.Version
	buf[StatusOffset] = h.Status
	binary.BigEndian.PutUint16(buf[ServiceIdOffset:], h.ServiceId)
	binary.BigEndian.PutUint16(buf[CmdOffset:], h.Cmd)
	binary.BigEndian.PutUint32(buf[SeqOffset:], h.Seq)
	return buf
}

func (h *Header) decode(data []byte) {
	h.Version = data[VersionOffset]
	h.Status = data[StatusOffset]
	h.ServiceId = binary.BigEndian.Uint16(data[ServiceIdOffset:CmdOffset])
	h.Cmd = binary.BigEndian.Uint16(data[CmdOffset:SeqOffset])
	h.Seq = binary.BigEndian.Uint32(data[SeqOffset:])
}

type Message struct {
	Header
	Body []byte
}

type imCodec struct {
	conn net.Conn
}

func NewImCodec(conn net.Conn) *imCodec {
	return &imCodec{
		conn: conn,
	}
}

func (cc *imCodec) readPackLength() (uint32, error) {
	data := make([]byte, PackSize)
	if n, err := cc.conn.Read(data); err != nil {
		return 0, err
	} else if n != PackSize {
		return 0, ErrPackSizeLength(n)
	} else {
		return binary.BigEndian.Uint32(data), nil
	}
}

func (cc *imCodec) Receive() (*Message, error) {
	if packLen, err := cc.readPackLength(); err != nil {
		return nil, err
	} else {
		data := make([]byte, packLen)
		n, err := cc.conn.Read(data)
		if err != nil {
			return nil, err
		} else if n != int(packLen) {
			return nil, ErrReadLength(n, int(packLen))
		}
		headerLen := binary.BigEndian.Uint16(data[:HeaderSize])
		msg := &Message{}
		msg.Header.decode(data[HeaderSize : HeaderSize+headerLen])
		msg.Body = data[HeaderSize+headerLen:]
		return msg, nil
	}
}

func (cc *imCodec) Write(msg *Message) error {
	packLen := HeaderSize + HeaderLength + len(msg.Body)
	var data []byte
	data = binary.BigEndian.AppendUint32(data, uint32(packLen))
	data = binary.BigEndian.AppendUint16(data, uint16(HeaderLength))
	data = append(data, msg.encode()...)
	data = append(data, msg.Body...)
	n, err := cc.conn.Write(data)
	if err != nil {
		return err
	}
	if n != len(data) {
		return ErrWriteLength(n, len(data))
	}
	return nil
}
