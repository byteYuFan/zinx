package znet

import "github.com/byteYuFan/zinx/zinterfance"

// BaseRouter 基础路由的实现
// 实现router时，先嵌入这个BaseRouter的基类，之后按需修改
type BaseRouter struct {
}

// PreHandle 在处理conn业务之前的钩子方法Hook
func (b *BaseRouter) PreHandle(request zinterfance.IRequest) {

}

// Handle 在处理conn业务中的钩子方法Hook
func (b *BaseRouter) Handle(request zinterfance.IRequest) {

}

// PostHandle 在处理conn业务后的钩子方法Hook
func (b *BaseRouter) PostHandle(request zinterfance.IRequest) {

}
