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

// request 处理
type RequestsMiddlerware func(*Requests)

// response 处理
type ResponseMiddlerware func(*Response)

type SpiderBase interface {
	allBaseMiddlerware(requMiddlers []RequestsMiddlerware, respMiddlers []ResponseMiddlerware)
	spiderBase() (string, int)
	transportBase() (*http.Client, *[]Response)
	pushRequestToSche() <-chan *Requests

	Parse()
}

type Spider struct {
	DefaultHeader http.Header
	Timeout       time.Duration

	Name        string
	CtxRoot     context.Context
	Client      http.Client
	ReqChan     chan *Requests
	Responses   []Response
	Middlewares []Middlerware

	ConcurrentHttpCount  int // 并发请求数
	ConcurrentParseCount int // 并发处理数

	requestsCallback []RequestsMiddlerware // 请求前
	responseCallback []ResponseMiddlerware // 响应后
	ParseFunc        ParseFunc             // 解析方法

	lock   sync.RWMutex
	status status
}

// 解析
type ParseFunc func(*Response)

func NewSpider(name string) *Spider {
	b := &Spider{
		Name:                 name,
		CtxRoot:              context.Background(),
		Client:               http.Client{},
		ConcurrentHttpCount:  16,
		ConcurrentParseCount: 16,
	}
	b.ReqChan = make(chan *Requests, 1)
	b.Responses = make([]Response, 0)
	b.Middlewares = make([]Middlerware, 0)
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
	//	case req := <-b.ReqChan:
	//		b.decorationRequests(req)
	//	}
	//}()
}

// 添加请求
func (b *Spider) AddRequest(req *Requests) {
	b.decorationRequests(req)
	b.ReqChan <- req
}

func (b *Spider) OnRequests(f RequestsMiddlerware) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if f == nil {
		return
	}
	b.requestsCallback = append(b.requestsCallback, f)
}

func (b *Spider) OnResponse(f ResponseMiddlerware) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if f == nil {
		return
	}
	b.responseCallback = append(b.responseCallback, f)
}

// 获取爬虫基本信息
func (b *Spider) spiderBase() (string, int) {
	return b.Name, b.ConcurrentHttpCount
}

// 获取爬虫数据
func (b *Spider) transportBase() (*http.Client, *[]Response) {
	return &b.Client, &b.Responses
}

// 推到调度器中
func (b *Spider) pushRequestToSche() <-chan *Requests {
	b.lock.Lock()
	defer b.lock.Unlock()
	return b.ReqChan
}

// 调度器来调用，加载基本中间件
func (b *Spider) allBaseMiddlerware(requMiddlers []RequestsMiddlerware, respMiddlers []ResponseMiddlerware) {
	b.requestsCallback = requMiddlers
	b.responseCallback = respMiddlers
}

// 自定义处理 request
func (b *Spider) decorationRequests(req *Requests) {
	if b.requestsCallback == nil {
		return
	}
	for _, requestsCallback := range b.requestsCallback {
		requestsCallback(req)
	}
	return
}

// 自定义处理 response
func (b *Spider) decorationResponse(resp *Response) {

	if b.responseCallback == nil {
		return
	}
	for _, responseCallback := range b.responseCallback {
		responseCallback(resp)
	}
	return
}
