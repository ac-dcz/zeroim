package libnet

import (
	"net"
	"testing"
)

func TestCodec(t *testing.T) {
	msg := &Message{
		Header: Header{
			Version:   1,
			Status:    1,
			ServiceId: 2,
			Cmd:       5,
			Seq:       6,
		},
		Body: []byte{'a', 'b', 'c'},
	}
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		t.Fatal(err)
	}
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
	t.Log(data.Header)
	t.Log(data.Body)
}
