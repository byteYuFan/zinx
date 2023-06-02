package znet

import (
	"errors"
	"fmt"
	"github.com/byteYuFan/zinx/utils"
	"github.com/byteYuFan/zinx/zinterfance"
	"io"
	"net"
	"sync"
)

// Connection 连接实体
type Connection struct {
	// 当前 TCPServer Connection 是属于那个Server的
	TCPServer zinterfance.IServer
	// Conn 当前连接的 socket TCP 套接字
	Conn *net.TCPConn
	// ConnID 当前连接的ID
	ConnID uint32
	// isClosed 当前连接的状态
	isClosed bool
	// ExitChan 告知当前连接以及退出的channel
	ExitChan chan bool
	// 无缓存管道，用于读、写 goroutine之间的通信
	msgChan chan []byte
	// MsgHandler 当前server的消息管理模块，用来绑定MsgID和对应处理业务API的关系
	MsgHandler zinterfance.IMsgHandle
	// property 连接属性集合
	property map[string]any
	// propertyLock 保护连接属性修改的锁
	propertyLock sync.RWMutex
}

// NewConnection 根据传入的参数返回一个连接实体
func NewConnection(server zinterfance.IServer, conn *net.TCPConn, connID uint32, handle zinterfance.IMsgHandle) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: handle,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte),
		TCPServer:  server,
		property:   make(map[string]any),
	}
	c.TCPServer.GetConnManager().Add(c)
	return c
}

// StartReader 连接的读业务
func (c *Connection) StartReader() {
	fmt.Println("[Reader] Goroutine is Running ....")
	defer fmt.Println("ConnID=", c.ConnID, "Reader is Exit.Remote Addr is", c.RemoteAddr().String())
	defer c.Stop()
	for {
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//count, err := c.Conn.Read(buf)
		//
		//if err != nil {
		//	// 当读取数据时发生错误
		//	if err == io.EOF {
		//		// 远程主机关闭连接，退出连接处理循环
		//		fmt.Println("远程主机关闭连接")
		//		break
		//	} else if opErr, ok := err.(*net.OpError); ok {
		//		if strings.Contains(opErr.Error(), "An existing connection was forcibly closed by the remote host") {
		//			// 远程主机关闭连接，退出连接处理循环
		//			fmt.Println("远程主机关闭连接")
		//			break
		//		}
		//	}
		//
		//	// 其他错误处理
		//	fmt.Println("读取数据时发生其他错误:", err)
		//	continue
		//}
		//创建一个拆包解包的对象
		dp := NewDataPack()
		// 读取客户端的Msg Head 二进制流 8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.Conn, headData); err != nil {
			fmt.Println("Read msg head err", err)
			break
		}
		// 拆包获取msgID和msgDataLen
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("Read msg data error", data)
				break
			}
		}
		msg.SetMsgData(data)
		// 根据Len在读取data
		//得到当前连接的request数据
		req := &Request{
			conn: c,
			msg:  msg,
		}
		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 开启了工作池
			// 交给连接池机制
			c.MsgHandler.SendMsgToTask(req)
		} else {
			go c.MsgHandler.DoMsgHandler(req)
		}
	}
}

// StartWriter 专门发送给客户端的模块
func (c *Connection) StartWriter() {
	defer fmt.Println(c.RemoteAddr().String(), "[Conn Writer exit!]")
	fmt.Println("[Writer Goroutine in Running...]")
	// 不断的循环阻塞等待chan的消息,如果有消息，回写给客户端
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error", err, "Conn Writer exit")
				return
			}
		case <-c.ExitChan:
			// 代表Reader退出，此时Writer也要退出
			return

		}
	}
}

// Start 启动连接 让当前的连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start()....ConnID=", c.ConnID)
	go c.StartReader()
	//TODO 启动从当前连接写数据的业务
	go c.StartWriter()
	c.TCPServer.CallOnConnStart(c)
}

// Stop 停止连接 结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop()...ConnID=", c.ConnID)
	// 如果当前连接关闭直接return
	if c.isClosed == true {
		return
	}
	c.TCPServer.CallOnConnStop(c)
	//关闭通道
	c.ExitChan <- true
	close(c.ExitChan)
	close(c.msgChan)
	c.isClosed = true
	c.TCPServer.GetConnManager().Remove(c)
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

// SendMsg 提供一个SendMsg方法 将我们要发送给客户端的数据，先进行封包，在发送
func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("connection closed when send msg")
	}
	// 将data进行封包
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgID, data))
	if err != nil {
		fmt.Println("Pack error msg id=", msgID)
		return err
	}
	// 将数据写给channel
	c.msgChan <- binaryMsg
	return nil
}

// SetProperty 设置连接属性
func (c *Connection) SetProperty(key string, value any) {
	// 加写锁
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	c.property[key] = value
}

// GetProperty 获取连接属性
func (c *Connection) GetProperty(key string) (any, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if value, ok := c.property[key]; ok {
		return value, nil
	}
	return nil, errors.New("property not found")
}

// RemoverProperty 移除连接属性
func (c *Connection) RemoverProperty(key string) {
	// 加写锁
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	delete(c.property, key)
}
