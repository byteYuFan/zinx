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
	// 当前server的消息管理模块，用来绑定MsgID和对应处理业务API的关系
	MsgHandler zinterfance.IMsgHandle
	// ConnManager 连接管理器
	ConnManager zinterfance.IConnManager

	// Hook 函数在创建连接之前调用
	OnConnStart func(conn zinterfance.IConnection)
	// Hook 函数在销毁前调用
	OnConnStop func(conn zinterfance.IConnection)
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
		// 开启消息队列及worker工作
		s.MsgHandler.StartWorkerPool()
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
			// 设置最大连接个数的判断，如果超过最大连接数量，那么关闭新的连接
			if s.ConnManager.Len() >= utils.GlobalObject.MaxConn {
				// TODO 给用户响应一个错误包
				conn.Close()
				continue
			}
			// 将处理新连接业务的方法和Conn进行绑定，得到连接模块
			dialConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++
			go dialConn.Start()
		}

	}()
}

// Stop 实现IServer接口
func (s *Server) Stop() {
	//TODO 将一些服务器的资源，状态或者一些开辟的信息进行回收
	fmt.Println("[STOP] Zinx server name:", s.Name)
	s.ConnManager.ClearConn()
}

// Server 实现IServer接口
func (s *Server) Server() {
	s.Start()
	// TODO 做一些额外的业务
	select {}
}

// AddRouter 实现增加AddRouter接口
func (s *Server) AddRouter(msgID uint32, router zinterfance.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router Successfully!!")
}

func (s *Server) GetConnManager() zinterfance.IConnManager {
	return s.ConnManager
}

// NewServer 提供一个初始化Server模块的方法
func NewServer(name string) zinterfance.IServer {
	return &Server{
		Name:        utils.GlobalObject.Name,
		IPVersion:   "tcp4",
		IP:          utils.GlobalObject.Host,
		Port:        utils.GlobalObject.TcpPort,
		MsgHandler:  NewMsgHandle(),
		ConnManager: NewConnManager(),
	}

}

// SetOnConnStart 注册钩子函数的方法
func (s *Server) SetOnConnStart(hookFun func(connection zinterfance.IConnection)) {
	s.OnConnStart = hookFun
}

// CallOnConnStart 调用钩子函数的方法
func (s *Server) CallOnConnStart(connection zinterfance.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("------->Call onConnStart()")
		s.OnConnStart(connection)
	}
}

// SetOnConnStop 注册钩子函数的方法
func (s *Server) SetOnConnStop(hookFun func(connection zinterfance.IConnection)) {
	s.OnConnStop = hookFun
}

// CallOnConnStop 调用钩子函数的方法
func (s *Server) CallOnConnStop(connection zinterfance.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("------->Call onConnStop()")
		s.OnConnStop(connection)
	}
}
