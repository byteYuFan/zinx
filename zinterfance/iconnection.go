package zinterfance

import "net"

// IConnection 连接层抽象接口
type IConnection interface {
	// Start 启动连接 让当前的连接准备开始工作
	Start()
	// Stop 停止连接 结束当前连接的工作
	Stop()
	// GetTCPConnection 夺取当前连接绑定的 socket connect
	GetTCPConnection() *net.TCPConn
	// GetConnID 获取当前连接ID
	GetConnID() uint32
	// RemoteAddr 获取远程客户端的TCP 状态 ip port
	RemoteAddr() net.Addr
	// Send 发送数据
	Send(data []byte) error
}

// HandleFun HandleFun 定义一个处理连接的函数
type HandleFun func(*net.TCPConn, []byte, int) error
