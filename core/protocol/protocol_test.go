package protocol

import (
	"net"
	"testing"
	"time"
)

func TestCodec(t *testing.T) {
	msg := &IMMessage{
		H: Header{
			Version:    1,
			StatusCode: 200,
			MsgType:    PrivateChatType,
			Seq:        0,
		},
		Body: &PrivateChatMessage{
			From:        1,
			To:          0,
			Ts:          uint64(time.Now().Unix()),
			ContentType: 1,
			Data:        []byte{'a', 'b', 'c'},
		},
	}
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		t.Fatal(err)
	}
	defer listen.Close()
	go func() {
		conn, err := net.Dial("tcp", ":8080")
		if err != nil {
			panic(err)
		}
		cc := NewImCodec(conn)
		if err := cc.Write(msg); err != nil {
			panic(err)
		}
	}()

	conn, err := listen.Accept()
	if err != nil {
		t.Fatal(err)
	}
	cc := NewImCodec(conn)
	data, err := cc.Receive()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Header: %#v \n", data.H)
	t.Logf("Body: %#v \n", data.Body)

	if data.H != msg.H {
		t.Fatal("data error")
	}
}
