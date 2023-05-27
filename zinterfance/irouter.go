package zinterfance

// IRouter 路由的抽象接口
// 路由的数据都是IRequest
type IRouter interface {
	// PreHandle 在处理conn业务之前的钩子方法Hook
	PreHandle(request IRequest)
	// Handle 在处理conn业务中的钩子方法Hook
	Handle(request IRequest)
	// PostHandle 在处理conn业务后的钩子方法Hook
	PostHandle(request IRequest)
}
