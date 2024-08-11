package protocol

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

/*
总长度
----|
header头长度
--|
header头
-|-|--|--|----|------...
版本号|状态码|消息类型|命令|seq|pb body体

header头长度=1字节版本号+1字节状态码+2字节消息类型+2字节命令+4字节seq
总长度=header头+header头长度+pb body体长度

----|--|-|-|--|--|----|body
总长度|header头长度|版本号|状态码|消息类型|命令|seq｜body
总长度=2+1+1+2+2+4+len(body)
header头长度=1+1+2+2+4
*/

type Protocol interface {
	NewCodec(conn net.Conn) Codec
}

type ImProtocol struct{}

func (ImProtocol) NewCodec(conn net.Conn) Codec {
	return NewImCodec(conn)
}

type Codec interface {
	Receive() (*IMMessage, error)
	Write(msg *IMMessage) error
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
	Close() error
}

type imCodec struct {
	conn net.Conn
}

func NewImCodec(conn net.Conn) Codec {
	return &imCodec{
		conn: conn,
	}
}

func (cc *imCodec) readPackLength() (uint32, error) {
	data := make([]byte, packsizeLen)
	if n, err := cc.conn.Read(data); err != nil {
		return 0, err
	} else if n != packsizeLen {
		return 0, ErrPackSizeLength(n)
	} else {
		return binary.BigEndian.Uint32(data), nil
	}
}

func (cc *imCodec) Receive() (*IMMessage, error) {
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
		msg := &IMMessage{}

		headerSize := binary.BigEndian.Uint16(data[:headersizeLen])
		if headerSize != uint16(headersize) {
			return nil, ErrHeaderInVaild(int(headerSize))
		}
		if err := msg.H.Decode(data[headersizeLen : headersizeLen+headersize]); err != nil {
			return nil, fmt.Errorf("Decode Header: %v", err)
		}

		switch msg.H.MsgType {
		case PrivateChat:
			chat := &PrivateChatMessage{}
			if err := chat.Decode(data[headersizeLen+headersize:]); err != nil {
				return nil, fmt.Errorf("Decode Body: %v", err)
			}
			msg.Body = chat
		case RoomChat:
			chat := &RoomChatMessage{}
			if err := chat.Decode(data[headersizeLen+headersize:]); err != nil {
				return nil, fmt.Errorf("Decode Body: %v", err)
			}
			msg.Body = chat
		}

		return msg, nil
	}
}

func (cc *imCodec) Write(msg *IMMessage) (err error) {
	var headerData, bodyData []byte
	if headerData, err = msg.H.Encode(); err != nil {
		return fmt.Errorf("Encode Header: %v", err)
	}
	if bodyData, err = msg.Body.Encode(); err != nil {
		return fmt.Errorf("Encode Body: %v", err)
	}
	var data []byte
	packSize := headersizeLen + len(headerData) + len(bodyData)
	data = binary.BigEndian.AppendUint32(data, uint32(packSize))
	data = binary.BigEndian.AppendUint16(data, uint16(len(headerData)))
	data = append(data, headerData...)
	data = append(data, bodyData...)

	n, err := cc.conn.Write(data)

	if err != nil {
		return err
	}

	if n != len(data) {
		return ErrWriteLength(n, len(data))
	}

	return nil
}

func (c *imCodec) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *imCodec) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

func (c *imCodec) Close() error {
	return c.conn.Close()
}