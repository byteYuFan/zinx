package znet

import (
	"errors"
	"fmt"
	"github.com/byteYuFan/zinx/zinterfance"
	"sync"
)

// ConnManager 连接管理实体
type ConnManager struct {
	// connections 管理的连接信息集合
	connections map[uint32]zinterfance.IConnection
	// connLock 读写保护连接集合的读写锁
	connLock sync.RWMutex
}

// NewConnManager 新建一个对象
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]zinterfance.IConnection),
	}
}

// Add 添加连接
func (cm *ConnManager) Add(conn zinterfance.IConnection) {
	// 保护共享资源锁 加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	// 将conn加入到Conn里面去
	cm.connections[conn.GetConnID()] = conn
	fmt.Println("ConnID=", conn.GetConnID(), "connection add to ConnManager successfully: conn num=", cm.Len())
}

// Remove 删除连接
func (cm *ConnManager) Remove(conn zinterfance.IConnection) {
	// 保护共享资源锁 加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	// 删除资源
	delete(cm.connections, conn.GetConnID())
	fmt.Println("ConnID=", conn.GetConnID(), "connection remove successfully!")
}

// Get 根据 connID获取连接
func (cm *ConnManager) Get(connID uint32) (zinterfance.IConnection, error) {
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()
	if conn, ok := cm.connections[connID]; ok {
		return conn, nil
	}
	return nil, errors.New("connection not found")
}

// Len 获取当前连接总数
func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

// ClearConn 清除连接
func (cm *ConnManager) ClearConn() {
	// 保护共享资源锁 加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	// 删除
	for connID, conn := range cm.connections {
		// 删除
		conn.Stop()
		delete(cm.connections, connID)
	}
	fmt.Println("Clear All Connections!!! nums=", cm.Len())
}
