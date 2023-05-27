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

// PreHandle 重写
func (p *PingRouter) PreHandle(request zinterfance.IRequest) {
	fmt.Println("Call Router PreHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
	if err != nil {
		fmt.Println("call back before ping err", err)
	}
}

// Handle 重写
func (p *PingRouter) Handle(request zinterfance.IRequest) {
	fmt.Println("Call Router Handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...\n"))
	if err != nil {
		fmt.Println("call back ping err", err)
	}
}

// PostHandle 重写
func (p *PingRouter) PostHandle(request zinterfance.IRequest) {
	fmt.Println("Call Router PostHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping ...\n"))
	if err != nil {
		fmt.Println("call back after ping err", err)
	}
}
func main() {
	s := znet.NewServer("zinx-v.0.3")
	s.AddRouter(new(PingRouter))
	s.Server()
}
