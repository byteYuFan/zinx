package zinterfance

// IMsgHandle 消息管理抽象层
type IMsgHandle interface {
	// DoMsgHandler 执行对应的Router消息处理方法
	DoMsgHandler(request IRequest)
	// AddRouter 为消息添加处理逻辑
	AddRouter(msgID uint32, router IRouter)
	// StartWorkerPool 启动worker工作池
	StartWorkerPool()
	// SendMsgToTask 将消息发送给消息队列去处理
	SendMsgToTask(request IRequest)
}
