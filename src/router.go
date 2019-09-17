package tgo

// 路由
type Router interface {
	// 匹配Handler
	//MatchHandler(ctx Context)  Handler
	Handle(ctx Context)
}

type RouterMatchHandlerFnc func(ctx Context,handlerMap map[interface{}]Handler) Handler
type DefaultRouter struct {
	matchHandlerFnc RouterMatchHandlerFnc
	handlerMap map[interface{}]Handler
}

func NewDefaultRouter(matchHandlerFnc RouterMatchHandlerFnc) *DefaultRouter  {
	return &DefaultRouter{matchHandlerFnc:matchHandlerFnc,handlerMap: map[interface{}]Handler{}}
}

func (d *DefaultRouter)  matchHandler(ctx Context)  Handler {
	return d.matchHandlerFnc(ctx,d.handlerMap)
}

func (d *DefaultRouter) Handle(ctx Context) {
	handle := d.matchHandler(ctx)
	if handle!=nil {
		 handle(ctx)
	}
}

func (d *DefaultRouter) Route(path interface{},handler Handler)  {
   d.handlerMap[path] = handler
}