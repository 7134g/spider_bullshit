package connect

import (
	"bytes"
	"container/list"
	"io"
	"net/http"
	"sync"
)

type Requests struct {
	SpiderName string
	Request    *http.Request
}

func (b *Spider) NewRequests(method, u string, body []byte) (*Requests, error) {
	var reqBody io.Reader
	if body == nil {
		reqBody = nil
	} else {
		reqBody = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(b.CtxRoot, method, u, reqBody)
	if err != nil {
		return nil, err
	}

	requests := &Requests{}
	requests.SpiderName = b.Name
	requests.Request = req
	req.Header = b.DefaultHeader

	return requests, nil
}

// 请求队列
type RequQueue interface {
	Push(v *Requests)
	Get() *Requests
	Queue() interface{}
}

// 限制长度的请求队列
type RequCh struct {
	ch chan *Requests
}

func NewRequCh(size int) *RequCh {
	return &RequCh{ch: make(chan *Requests, size)}
}

// 尾插
func (l *RequCh) Push(v *Requests) {
	if v != nil {
		l.ch <- v
	}
}

// 从头取
func (l *RequCh) Get() *Requests {
	return <-l.ch
}

func (l *RequCh) Queue() interface{} {
	return l.ch
}

// 无限制请求队列
type RequList struct {
	list *list.List
	mux  sync.Mutex
}

func NewRequList() *RequList {
	return &RequList{list: list.New(), mux: sync.Mutex{}}
}

// 尾插
func (l *RequList) Push(v *Requests) {
	l.mux.Lock()
	l.list.PushFront(v)
	l.mux.Unlock()
}

// 从头取
func (l *RequList) Get() *Requests {
	l.mux.Lock()
	defer l.mux.Unlock()
	var inter interface{}

	if l.Len() > 0 {
		inter = l.list.Remove(l.list.Back())
		return inter.(*Requests)
	} else {
		return nil
	}
}

func (l *RequList) Queue() interface{} {
	return l.list
}

// 移除
func (l *RequList) Remove(e *list.Element) {
	l.list.Remove(e)
}

// 获取长度
func (l *RequList) Len() int {
	length := l.list.Len()
	return length
}
