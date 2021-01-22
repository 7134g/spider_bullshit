package connect

// request 处理
type RequestsMiddleware func(*Requests) error

// response 处理
type ResponseMiddleware func(*Response) error

// 解析
type ParseFunc func(*Response) interface{}

// 中间件
type Middleware struct {
	Name           string
	RequMiddleware RequestsMiddleware
	RespMiddleware ResponseMiddleware
}

func NewMiddleware(name string) Middleware {
	return Middleware{
		Name: name,
	}
}

func (m *Middleware) ProcessRequests(req *Requests) error {
	if m.RequMiddleware != nil && req != nil {
		err := m.RequMiddleware(req)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (m *Middleware) ProcessResponse(res *Response) error {
	if m.RespMiddleware != nil && res != nil {
		err := m.RespMiddleware(res)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}
