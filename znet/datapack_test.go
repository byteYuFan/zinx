package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

// TestDataPack 测试封装
func TestDataPack(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		t.Log(err)
		panic(err)
	}
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				t.Log(err)
				continue
			}
			go func(conn net.Conn) {
				dp := NewDataPack()
				for {
					headData := make([]byte, dp.GetHeadLen())
					_, err = io.ReadFull(conn, headData)
					if err != nil {
						t.Log(err)
						return
					}
					msgHead, err := dp.Unpack(headData)
					if err != nil {
						t.Log(err)
						return
					}

					if msgHead.GetMsgLen() > 0 {
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							t.Log(err)
							return
						}
						fmt.Println("receive message-------------->", msg.ID, msg.DataLen, string(msg.Data))
					}
				}
			}(conn)
		}
	}()
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client err", err)
		return
	}
	dp := NewDataPack()
	msg1 := &Message{
		ID:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println(err)
		return
	}
	msg2 := &Message{
		ID:      2,
		DataLen: 7,
		Data:    []byte{'h', 'e', 'l', 'l', 'o', 'm', 'm'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println(err)
		return
	}
	sendData1 = append(sendData1, sendData2...)
	_, err = conn.Write(sendData1)
	if err != nil {
		fmt.Println(err)
		return
	}
	select {}

}
