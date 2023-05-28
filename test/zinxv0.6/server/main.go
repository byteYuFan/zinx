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
	fmt.Println("Call ping Router Handle")
	// 客户端发送来消息后，在进行回写
	fmt.Println("Receive from client msgID=", request.GetID(),
		"data=", string(request.GetData()))
	request.GetData()
	err := request.GetConnection().SendMsg(200, []byte("ping...ping..."))
	if err != nil {
		fmt.Println(err)
	}
}

// HelloRouter PingRouter测试
type HelloRouter struct {
	znet.BaseRouter
}

// Handle 重写
func (p *HelloRouter) Handle(request zinterfance.IRequest) {
	fmt.Println("Call hello Router Handle")
	// 客户端发送来消息后，在进行回写
	fmt.Println("Receive from client msgID=", request.GetID(),
		"data=", string(request.GetData()))
	request.GetData()
	err := request.GetConnection().SendMsg(201, []byte("hello...hello...welcome to zinx"))
	if err != nil {
		fmt.Println(err)
	}
}

// Start 连接之后的函数
func Start(conn zinterfance.IConnection) {
	fmt.Println("[Start........................]")
	conn.SendMsg(202, []byte("do start"))
}
func main() {
	s := znet.NewServer("zinx-v.0.5")
	s.SetOnConnStart(Start)
	s.AddRouter(0, new(PingRouter))
	s.AddRouter(1, new(HelloRouter))

	s.Server()
}
