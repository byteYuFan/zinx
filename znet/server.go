package znet

import (
	"fmt"
	"github.com/byteYuFan/zinx/zinterfance"
	"net"
)

// Server 定义一个Server的服务模块
type Server struct {
	// Name 服务器的名 称
	Name string
	// IPVersion 服务器绑定的ip版本
	IPVersion string
	// IP 服务器监听端口号
	IP string
	// Port 服务器监听的端口号
	Port int
}

// Start 实现IServer接口
func (s *Server) Start() {
	fmt.Printf("[Start] Server Listenner at IP:%s,Port %d is  starting\n", s.IP, s.Port)
	go func() {
		// 获取一个TCP Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error", err)
			return
		}
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err:", err)
			return
		}
		fmt.Println("start zinx server successfully", s.Name, "success,Listening...")
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}
			go func() {
				for {
					buf := make([]byte, 512)
					count, err := conn.Read(buf)
					if err != nil {
						fmt.Println("receive buf err", err)
						continue
					}
					if _, err := conn.Write(buf[:count]); err != nil {
						fmt.Println("write back buf err", err)
						continue
					}
				}
			}()
		}

	}()
}

// Stop 实现IServer接口
func (s *Server) Stop() {
	//TODO 将一些服务器的资源，状态或者一些开辟的信息进行回收
}

// Server 实现IServer接口
func (s *Server) Server() {
	s.Start()
	// TODO 做一些额外的业务
	select {}
}

// NewServer 提供一个初始化Server模块的方法
func NewServer(name string) zinterfance.IServer {
	return &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}

}
