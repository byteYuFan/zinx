package main

import (
	"fmt"
	"github.com/byteYuFan/zinx/znet"
	"io"
	"net"
	"time"
)

func main() {
	fmt.Println("client start...")
	time.Sleep(1 * time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err", err)
		return
	}
	for {
		dp := znet.NewDataPack()
		binary, err := dp.Pack(znet.NewMsgPackage(0, []byte("zinxv0.6 client test message.")))
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = conn.Write(binary)
		if err != nil {
			fmt.Println(err)
			return
		}
		// 先读取流中的head部分
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head error", err)
			break
		}
		msgHead, _ := dp.Unpack(binaryHead)
		if msgHead.GetMsgLen() > 0 {
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data error,", err)
				return
			}
			fmt.Println("----->receive server msg id=", msg.ID, "len=", msg.DataLen, "data=", string(msg.Data))
		}
		time.Sleep(1 * time.Second)
	}
}
