package znet

import (
	"fmt"
	"github.com/byteYuFan/zinx/utils"
	"github.com/byteYuFan/zinx/zinterfance"
	"strconv"
)

// MsgHandle 消息处理模块实现层
type MsgHandle struct {
	// Apis 消息ID 和 router 对应关系的集合
	Apis map[uint32]zinterfance.IRouter
	// 消息队列 负责worker读取任务的消息队列
	TaskQueue []chan zinterfance.IRequest
	// 负责业务Worker池的数量
	WorkerPoolSize uint32
}

// NewMsgHandle  创建方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]zinterfance.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize, // 全局配置中获取
		TaskQueue:      make([]chan zinterfance.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

// DoMsgHandler 执行对应的Router消息处理方法
func (mh *MsgHandle) DoMsgHandler(request zinterfance.IRequest) {
	// 获取id ,根据id进行查找
	msgID := request.GetID()
	handler, ok := mh.Apis[msgID]
	if !ok {
		fmt.Println("api msgID=", request.GetID(), "is NOT FOUND! YOU NEED REGISTER IT!")
		return
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// AddRouter 为消息添加处理逻辑
func (mh *MsgHandle) AddRouter(msgID uint32, router zinterfance.IRouter) {
	// 当前的msg绑定的api是否已经存在，如果存在则返回
	if _, ok := mh.Apis[msgID]; ok {
		// 表明api已经注册
		panic("repeat api,msgID=" + strconv.Itoa(int(msgID)))
		return
	}
	// 添加api
	mh.Apis[msgID] = router
	fmt.Println("Add api msgID", msgID, " successfully!")
}

// StartWorkerPool 启动一个worker工作池,开启工作池的动作只能发生一次
func (mh *MsgHandle) StartWorkerPool() {
	// 根据workerPoolSize 分别开启worker，每个worker用go承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 给当前的worker对应的channel开辟空间
		mh.TaskQueue[i] = make(chan zinterfance.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		// 启动当前的worker
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}

}

// StartOneWorker 启动一个Worker的工作流程
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan zinterfance.IRequest) {
	fmt.Println("Worker ID =", workerID, " is started")
	// 不断的阻塞等待对应消息队列的消息
	for {
		select {
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

// SendMsgToTask 将消息交给taskQueue
func (mh *MsgHandle) SendMsgToTask(request zinterfance.IRequest) {
	// 将消息平均分配给不同的worker
	// 根据客户端建立的ConnID来进行分配
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	// 将消息发送给对应worker的TaskQueue
	fmt.Println("Add ConnID=", request.GetConnection().GetConnID(),
		"request MsgID=", request.GetID(), " to WorkerID=", workerID)
	mh.TaskQueue[workerID] <- request
}
