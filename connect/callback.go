package connect

// request 处理
type RequestsMiddlerware func(*Requests, *Response) error

// response 处理
type ResponseMiddlerware func(*Requests, *Response) error

// 解析
type ParseFunc func(*Response)
