package main

import (
	"fmt"
	"github.com/byteYuFan/zinx/zinterfance"
	"github.com/byteYuFan/zinx/znet"
)

// PingRouter PingRouter测试
type PingRouter struct {
	znet.BaseRouter
}

// Handle 重写
func (p *PingRouter) Handle(request zinterfance.IRequest) {
	fmt.Println("Call Router Handle")
	// 客户端发送来消息后，在进行回写
	fmt.Println("Receive from client msgID=", request.GetID(),
		"data=", string(request.GetData()))
	request.GetData()
	err := request.GetConnection().SendMsg(1, []byte("ping...ping..."))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	s := znet.NewServer("zinx-v.0.5")
	//s.AddRouter(new(PingRouter))
	s.Server()
}
