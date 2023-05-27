package znet

import (
	"fmt"
	"github.com/byteYuFan/zinx/utils"
	"github.com/byteYuFan/zinx/zinterfance"
	"io"
	"net"
	"strings"
)

// Connection 连接实体
type Connection struct {
	// Conn 当前连接的 socket TCP 套接字
	Conn *net.TCPConn
	// ConnID 当前连接的ID
	ConnID uint32
	// isClosed 当前连接的状态
	isClosed bool
	// ExitChan 告知当前连接以及退出的channel
	ExitChan chan bool

	// Router 该连接处理的方法
	Router zinterfance.IRouter
}

// NewConnection 根据传入的参数返回一个连接实体
func NewConnection(conn *net.TCPConn, connID uint32, router zinterfance.IRouter) *Connection {
	return &Connection{
		Conn:     conn,
		ConnID:   connID,
		Router:   router,
		isClosed: false,
		ExitChan: make(chan bool, 1),
	}
}

// StartReader 连接的读业务
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is Running ....")
	defer fmt.Println("ConnID=", c.ConnID, "Reader is Exit.Remote Addr is", c.RemoteAddr().String())
	defer c.Stop()
	for {
		buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		count, err := c.Conn.Read(buf)

		if err != nil {
			// 当读取数据时发生错误
			if err == io.EOF {
				// 远程主机关闭连接，退出连接处理循环
				fmt.Println("远程主机关闭连接")
				break
			} else if opErr, ok := err.(*net.OpError); ok {
				if strings.Contains(opErr.Error(), "An existing connection was forcibly closed by the remote host") {
					// 远程主机关闭连接，退出连接处理循环
					fmt.Println("远程主机关闭连接")
					break
				}
			}

			// 其他错误处理
			fmt.Println("读取数据时发生其他错误:", err)
			continue
		}
		//得到当前连接的request数据
		req := &Request{
			conn: c,
			data: buf[:count],
		}
		// 从路由中找到注册绑定的conn对应的router调用
		go func(request zinterfance.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(req)

	}
}

// Start 启动连接 让当前的连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start()....ConnID=", c.ConnID)
	go c.StartReader()
	//TODO 启动从当前连接写数据的业务
}

// Stop 停止连接 结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop()...ConnID=", c.ConnID)
	// 如果当前连接关闭直接return
	if c.isClosed == true {
		return
	}
	//关闭通道
	close(c.ExitChan)
	c.isClosed = true
	//调用关闭socket连接
	err := c.Conn.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

}

// GetTCPConnection 或取当前连接绑定的 socket connect
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// GetConnID 获取当前连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// RemoteAddr 获取远程客户端的TCP 状态 ip port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// Send 发送数据
func (c *Connection) Send(data []byte) error {
	_, err := c.Conn.Write(data)
	if err != nil {
		return err
	}
	return nil
}
