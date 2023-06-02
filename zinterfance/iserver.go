package zinterfance

// IServer 定义一个服务器接口
type IServer interface {
	// Start 启动服务
	Start()
	// Stop 停止服务
	Stop()
	//Server 运行服务
	Server()
	// AddRouter 给当前的服务注册一个路由方法
	AddRouter(msgID uint32, router IRouter)
	// GetConnManager 获取当前server的连接管理器
	GetConnManager() IConnManager
	// SetOnConnStart 注册钩子函数的方法
	SetOnConnStart(func(connection IConnection))
	// CallOnConnStart 调用钩子函数的方法
	CallOnConnStart(connection IConnection)
	// SetOnConnStop 注册钩子函数的方法
	SetOnConnStop(func(connection IConnection))
	// CallOnConnStop 调用钩子函数的方法
	CallOnConnStop(connection IConnection)
}
