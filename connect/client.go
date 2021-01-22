package connect

import (
	"context"
	"net/http"
	"sync"
	"time"
)

type status int

const (
	_ = status(iota)
	running
	stop
)

type SpiderBase interface {
	allBaseMiddlerware(requMiddlers []RequestsMiddlerware, respMiddlers []ResponseMiddlerware)
	spiderBase() (string, int)
	transportBase() (*http.Client, RequQueue, RespQueue)
	pushRequestToSche() RequQueue

	getParse() ([]ResponseMiddlerware, ParseFunc)
}

type Spider struct {
	DefaultHeader http.Header
	Timeout       time.Duration

	Name      string
	CtxRoot   context.Context
	Client    http.Client
	RequQueue RequQueue
	RespList  RespQueue

	ConcurrentHttpCount  int // 并发请求数
	ConcurrentParseCount int // 并发处理数

	requestsMiddlewares []RequestsMiddlerware // 请求前
	responseMiddlewares []ResponseMiddlerware // 响应后
	ParseFunc           ParseFunc             // 解析方法

	lock   sync.RWMutex
	status status
}

func NewSpider(name string) *Spider {
	b := &Spider{
		Name:                 name,
		CtxRoot:              context.Background(),
		Client:               http.Client{},
		ConcurrentHttpCount:  16,
		ConcurrentParseCount: 16,
	}
	b.RequQueue = NewRequList()
	b.RespList = NewRespList()
	b.Client.Timeout = b.Timeout
	return b
}

func (b *Spider) Run() {
	b.lock.Lock()
	if b.status == running {
		return
	}
	b.status = running
	b.lock.Unlock()

	//go func() {
	//	select {
	//	case req := <-b.RequQueue:
	//		b.decorationRequests(req)
	//	}
	//}()
}

// 添加请求
func (b *Spider) AddRequest(req *Requests) {
	//b.decorationRequests(req)
	b.RequQueue.Push(req)
}

func (b *Spider) OnRequests(f RequestsMiddlerware) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if f == nil {
		return
	}
	b.requestsMiddlewares = append(b.requestsMiddlewares, f)
}

func (b *Spider) OnResponse(f ResponseMiddlerware) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if f == nil {
		return
	}
	b.responseMiddlewares = append(b.responseMiddlewares, f)
}

func (b *Spider) getParse() ([]ResponseMiddlerware, ParseFunc) {
	// 解析 response
	return b.responseMiddlewares, b.ParseFunc
}

// 获取爬虫基本信息
func (b *Spider) spiderBase() (string, int) {
	return b.Name, b.ConcurrentHttpCount
}

// 获取爬虫数据
func (b *Spider) transportBase() (*http.Client, RequQueue, RespQueue) {
	return &b.Client, b.RequQueue, b.RespList
}

// 推到调度器中
func (b *Spider) pushRequestToSche() RequQueue {
	b.lock.Lock()
	defer b.lock.Unlock()
	return b.RequQueue
}

// 调度器来调用，加载基本中间件
func (b *Spider) allBaseMiddlerware(requMs []RequestsMiddlerware, respMs []ResponseMiddlerware) {
	b.requestsMiddlewares = requMs
	b.responseMiddlewares = respMs
}

// 自定义处理 request
//func (b *Spider) decorationRequests(req *Requests) {
//	if b.requestsMiddlewares == nil {
//		return
//	}
//	for _, requestsM := range b.requestsMiddlewares {
//		requestsM(req, nil)
//	}
//	return
//}

// 自定义处理 response
//func (b *Spider) decorationResponse(resp *Response) {
//
//	if b.responseMiddlewares == nil {
//		return
//	}
//	for _, responseM := range b.responseMiddlewares {
//		responseM(nil, resp)
//	}
//	return
//}

func (b *Spider) GET(u string) (*Requests, error) {
	return b.NewRequests("GET", u, nil)
}

func (b *Spider) POST(u string, body []byte) (*Requests, error) {
	return b.NewRequests("POST", u, body)
}

func (b *Spider) HEAD(u string) (*Requests, error) {
	return b.NewRequests("HEAD", u, nil)
}

func (b *Spider) OPTIONS(u string) (*Requests, error) {
	return b.NewRequests("OPTIONS", u, nil)
}
