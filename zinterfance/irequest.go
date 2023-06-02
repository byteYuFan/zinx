package zinterfance

// IRequest 接口
// 客户端请求的连接和数据包装到了一起
type IRequest interface {
	// GetConnection 得到当前连接
	GetConnection() IConnection
	// GetData 得到请求数据
	GetData() []byte
	// GetID 获取请求的ID
	GetID() uint32
}
