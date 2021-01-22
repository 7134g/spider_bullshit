package connect

import (
	"bullshit/common/logs"
	"bullshit/pool"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"
)

type Scheduler struct {
	lock       sync.RWMutex
	SpiderList []SpiderBase  // 所有爬虫类型
	Groups     []*pool.Group // 所有工作组

	status                  status // 调度器状态
	baseRequestsMiddlerware []RequestsMiddlerware
	baseResponseMiddlerware []ResponseMiddlerware
}

// 获取一个调度器
func NewScheduler() *Scheduler {
	return &Scheduler{
		SpiderList: make([]SpiderBase, 0),
		Groups:     make([]*pool.Group, 0),
	}
}

// 添加基础共有中间件
func (s *Scheduler) AddBaseMiddlerware(requM RequestsMiddlerware, respM ResponseMiddlerware) {
	if requM != nil {
		s.baseRequestsMiddlerware = append(s.baseRequestsMiddlerware, requM)
	}
	if respM != nil {
		s.baseResponseMiddlerware = append(s.baseResponseMiddlerware, respM)
	}
}

// 添加新的爬虫
func (s *Scheduler) AddSpider(spider SpiderBase) {
	s.lock.Lock()
	defer s.lock.Unlock()
	// 装载共有模块
	spider.allBaseMiddlerware(s.baseRequestsMiddlerware, s.baseResponseMiddlerware)
	s.SpiderList = append(s.SpiderList, spider)
}

// 开始运行程序
func (s *Scheduler) Excute() {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.status == running {
		return
	}
	s.status = running

	// 为每一个爬虫创建工作组, 添加任务执行到工作组
	for _, spider := range s.SpiderList {
		name, count := spider.spiderBase()
		group := pool.AddNewGroup(name, count)
		s.Groups = append(s.Groups, group)
		go monitor(group, spider) // 处理新的请求
		go process(group, spider) // 处理响应的请求
	}
	s.status = stop

}

func SpiderStatus(g *pool.Group) bool {
	return g.Status
}

// 监控爬虫，处理该爬虫新的请求
func monitor(g *pool.Group, s SpiderBase) {
	// 发请求
	name, _ := s.spiderBase()
	client, requQueue, respQueue := s.transportBase()
	middlerwares, _ := s.getParse()

	for {
		if !SpiderStatus(g) {
			logs.Logger.Info(name + " is close")
			return
		}

		// 获取 resp
		requ := requQueue.Get()

		// 处理自定义中间件
		for _, middlerware := range middlerwares {
			err := middlerware(requ, nil)
			if err != nil {
				continue
			}
		}

		// 值为nil，若中间件不处理
		if requ == nil {
			time.Sleep(NO_RESPONSE_DATA_SLEEP_TIME)
		}

		// 发送请求
		t := &pool.Task{}
		t.TaskFunc = dohttp
		t.Param = []interface{}{requ, client, respQueue}
		err := g.Pool.Submit(t)
		if err != nil {
			logs.Logger.Error("g.Pool.Submit", zap.Error(err))
		}
	}

}

// 处理该爬虫新的响应数据
func process(g *pool.Group, s SpiderBase) {
	// 处理请求
	name, _ := s.spiderBase()
	_, _, respList := s.transportBase()
	middlerwares, parseFunc := s.getParse()
	for {
		if !SpiderStatus(g) {
			logs.Logger.Info(name + " is close")
			return
		}

		// 获取 resp
		resp := respList.Get()

		// 处理自定义中间件
		for _, middlerware := range middlerwares {
			err := middlerware(nil, resp)
			if err != nil {
				continue
			}
		}

		// 值为nil，若中间件不处理
		if resp == nil {
			time.Sleep(NO_RESPONSE_DATA_SLEEP_TIME)
		}

		// 开始解析
		t := &pool.Task{}
		t.TaskFunc = doparse
		t.Param = []interface{}{parseFunc, resp}
		err := g.Pool.Submit(t)
		if err != nil {
			logs.Logger.Error("g.Pool.Submit", zap.Error(err))
		}
	}
}

func dohttp(params []interface{}) {
	req := params[0].(*Requests)
	client := params[1].(*http.Client)
	respList := params[2].(*RespList)

	resp, err := client.Do(req.Request)
	if err != nil {
		logs.Logger.Error("http error", zap.Error(err))
		return
	}
	if resp == nil {
		logs.Logger.Error("response is nil")
		return
	}

	defer func() { _ = resp.Body.Close() }()

	response, err := NewResponse(req.SpiderName, resp)
	if err != nil {
		logs.Logger.Error("read body error", zap.Error(err))
		return
	}

	respList.Push(response)

}

func doparse(params []interface{}) {
	f := params[0].(ParseFunc)
	r := params[1].(*Response)

	f(r)
}
