package znet

import (
	"fmt"
	"github.com/byteYuFan/zinx/utils"
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
	// 当前的Server添加一个router
	Router zinterfance.IRouter
}

// Start 实现IServer接口
func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name :%s, Listener at IP :%s, Port:%d is starting\n",
		utils.GlobalObject.Name,
		utils.GlobalObject.Host,
		utils.GlobalObject.TcpPort,
	)
	fmt.Printf("[Zinx] Version %s ,MaxConn:%d,MaxPackageSize:%d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPackageSize,
	)
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
		var cid uint32

		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}
			// 将处理新连接业务的方法和Conn进行绑定，得到连接模块
			dialConn := NewConnection(conn, cid, s.Router)
			cid++
			go dialConn.Start()
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

// AddRouter 实现增加AddRouter接口
func (s *Server) AddRouter(router zinterfance.IRouter) {
	s.Router = router
	fmt.Println("Add Router Successfully!!")
}

// NewServer 提供一个初始化Server模块的方法
func NewServer(name string) zinterfance.IServer {
	return &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		Router:    nil,
	}

}
