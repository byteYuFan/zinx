package zinterfance

import (
	"net"
)

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
	// SendMsg 发送数据
	SendMsg(id uint32, data []byte) error
	// SetProperty 设置连接属性
	SetProperty(key string, value any)
	// GetProperty 获取连接属性
	GetProperty(key string) (any, error)
	// RemoverProperty 移除连接属性
	RemoverProperty(key string)
}

// HandleFun HandleFun 定义一个处理连接的函数
type HandleFun func(*net.TCPConn, []byte, int) error
