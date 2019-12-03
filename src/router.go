package tgo

// Router 路由
type Router interface {
	// 匹配Handler
	//MatchHandler(ctx Context)  Handler
	Handle(ctx Context)
}

// RouterMatchHandlerFnc 路由匹配函数
type RouterMatchHandlerFnc func(ctx Context, handlerMap map[interface{}]Handler) Handler

// DefaultRouter 默认路由
type DefaultRouter struct {
	matchHandlerFnc RouterMatchHandlerFnc
	handlerMap      map[interface{}]Handler
}

// NewDefaultRouter 创建一个默认路由
func NewDefaultRouter(matchHandlerFnc RouterMatchHandlerFnc) *DefaultRouter {
	return &DefaultRouter{matchHandlerFnc: matchHandlerFnc, handlerMap: map[interface{}]Handler{}}
}

func (d *DefaultRouter) matchHandler(ctx Context) Handler {
	return d.matchHandlerFnc(ctx, d.handlerMap)
}

// Handle 处理请求
func (d *DefaultRouter) Handle(ctx Context) {
	handle := d.matchHandler(ctx)
	if handle != nil {
		handle(ctx)
	}
}

// Route 设置路由
func (d *DefaultRouter) Route(path interface{}, handler Handler) {
	d.handlerMap[path] = handler
}
