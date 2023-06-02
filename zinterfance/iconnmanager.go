package zinterfance

// IConnManager 连接管理模块:抽象层
type IConnManager interface {
	// Add 添加连接
	Add(conn IConnection)
	// Remove 删除连接
	Remove(conn IConnection)
	// Get 根据 connID获取连接
	Get(connID uint32) (IConnection, error)
	// Len 获取当前连接总数
	Len() int
	// ClearConn 清除连接
	ClearConn()
}
