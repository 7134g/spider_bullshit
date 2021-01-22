package connect

import (
	"container/list"
	"io/ioutil"
	"net/http"
	"sync"
)

type Response struct {
	SpiderName string
	Resp       *http.Response
	Body       []byte
}

func NewResponse(name string, resp *http.Response) (*Response, error) {
	var err error
	response := &Response{}
	response.SpiderName = name
	response.Resp = resp
	response.Body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return response, err
}

// 响应队列
type RespQueue interface {
	Push(v *Response)
	Get() *Response
	Queue() interface{}
}

// 限制长度的结果队列
type RespCh struct {
	ch chan *Response
}

func NewRespCh(size int) *RespCh {
	return &RespCh{ch: make(chan *Response, size)}
}

// 尾插
func (l *RespCh) Push(v *Response) {
	if v != nil {
		l.ch <- v
	}
}

// 从头取
func (l *RespCh) Get() *Response {
	return <-l.ch
}

func (l *RespCh) Queue() interface{} {
	return l.ch
}

// 无限制结果队列
type RespList struct {
	list *list.List
	mux  sync.Mutex
}

func NewRespList() *RespList {
	return &RespList{list: list.New(), mux: sync.Mutex{}}
}

// 尾插
func (l *RespList) Push(v *Response) {
	l.mux.Lock()
	l.list.PushFront(v)
	l.mux.Unlock()
}

// 从头取
func (l *RespList) Get() *Response {
	l.mux.Lock()
	defer l.mux.Unlock()
	var inter interface{}

	if l.Len() > 0 {
		inter = l.list.Remove(l.list.Back())
		return inter.(*Response)
	} else {
		return nil
	}
}

func (l *RespList) Queue() interface{} {
	return l.list
}

// 移除
func (l *RespList) Remove(e *list.Element) {
	l.list.Remove(e)
}

// 获取长度
func (l *RespList) Len() int {
	length := l.list.Len()
	return length
}
