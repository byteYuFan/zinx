package znet

import "github.com/byteYuFan/zinx/zinterfance"

// Request 请求实体
type Request struct {
	// conn 已经和客户端建立好的连接
	conn zinterfance.IConnection
	// data 客户端请求好的数据
	data []byte
}

// GetConnection 得到当前连接
func (r *Request) GetConnection() zinterfance.IConnection {
	return r.conn
}

// GetData 得到请求数据
func (r *Request) GetData() []byte {
	return r.data
}
