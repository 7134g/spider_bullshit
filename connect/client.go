package connect

import (
	"context"
	"fmt"
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
	//allBaseMiddlerware(requMiddlers []RequestsMiddleware, respMiddlers []ResponseMiddleware)
	spiderBase() (string, int)
	transportBase() (*http.Client, RequQueue, RespQueue)
	pushRequestToSche() RequQueue

	loadMiddlerwares(ms []Middleware)
	getMiddlerwares() []Middleware
	getParse() ParseFunc
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

	Middlewares []Middleware // 中间件
	ParseFunc   ParseFunc    // 解析方法

	lock   sync.RWMutex // 同步操作
	status status       // 爬虫状态
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

func (b *Spider) OnRequests(f RequestsMiddleware) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if f == nil {
		return
	}
	mName := fmt.Sprintf("%s-Requests", b.Name)
	m := NewMiddleware(mName)
	m.RequMiddleware = f
	b.Middlewares = append(b.Middlewares, m)
	//b.requestsMiddlewares = append(b.requestsMiddlewares, f)
}

func (b *Spider) OnResponse(f ResponseMiddleware) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if f == nil {
		return
	}
	mName := fmt.Sprintf("%s-Response", b.Name)
	m := NewMiddleware(mName)
	m.RespMiddleware = f
	b.Middlewares = append(b.Middlewares, m)
	//b.responseMiddlewares = append(b.responseMiddlewares, f)
}

func (b *Spider) AddMiddleware(m Middleware) {
	b.Middlewares = append(b.Middlewares, m)
}

func (b *Spider) getParse() ParseFunc {
	// 解析 response
	return b.ParseFunc
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

func (b *Spider) loadMiddlerwares(ms []Middleware) {
	b.Middlewares = ms
}

func (b *Spider) getMiddlerwares() []Middleware {
	return b.Middlewares
}

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
