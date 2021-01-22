package connect

import (
	"bullshit/common/logs"
	"bullshit/pipeline"
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

	Middlewares []Middleware
	Pipeline    pipeline.Pipeline // 管道数据
	status      status            // 调度器状态
}

// 获取一个调度器
func NewScheduler() *Scheduler {
	return &Scheduler{
		SpiderList: make([]SpiderBase, 0),
		Groups:     make([]*pool.Group, 0),
	}
}

// 添加基础共有中间件
func (s *Scheduler) AddBaseMiddlerware(m Middleware) {
	if m.RespMiddleware != nil || m.RequMiddleware != nil {
		s.Middlewares = append(s.Middlewares, m)
	}
}

// 添加新的爬虫
func (s *Scheduler) AddSpider(spider SpiderBase) {
	s.lock.Lock()
	defer s.lock.Unlock()
	// 装载共有模块
	spider.loadMiddlerwares(s.Middlewares)
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
		httpGroup := pool.AddNewGroup(name, count)
		s.Groups = append(s.Groups, httpGroup)
		go monitor(httpGroup, spider)             // 处理新的请求
		go process(httpGroup, spider, s.Pipeline) // 处理响应的请求

		StorageGroup := pool.AddNewGroup(name, count)
		go storage(StorageGroup, s.Pipeline)
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
	middlerwares := s.getMiddlerwares()

	for {
		if !SpiderStatus(g) {
			logs.Logger.Info(name + " is close")
			return
		}

		// 获取 resp
		requ := requQueue.Get()

		// 处理自定义中间件
		for _, middlerware := range middlerwares {
			if middlerware.RequMiddleware == nil {
				continue
			}
			err := middlerware.ProcessRequests(requ)
			if err != nil {
				logs.Logger.Error(g.Name, zap.Error(err))
				return
			}
		}

		// 值为nil，若中间件不处理
		if requ == nil {
			time.Sleep(NO_REQUESTS_DATA_SLEEP_TIME)
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

	// todo 满足某些条件关闭 group 和 spider
	//g.Status = false
}

// 处理该爬虫新的响应数据
func process(g *pool.Group, s SpiderBase, p pipeline.Pipeline) {
	// 处理请求
	name, _ := s.spiderBase()
	_, _, respList := s.transportBase()
	parseFunc := s.getParse()
	middlerwares := s.getMiddlerwares()
	for {
		if !SpiderStatus(g) {
			logs.Logger.Info(name + " is close")
			return
		}

		// 获取 resp
		resp := respList.Get()

		// 处理自定义中间件
		for _, middlerware := range middlerwares {
			if middlerware.RespMiddleware == nil {
				continue
			}
			err := middlerware.ProcessResponse(resp)
			if err != nil {
				logs.Logger.Error(g.Name, zap.Error(err))
				return
			}
		}

		// 值为nil，若中间件不处理
		if resp == nil {
			time.Sleep(NO_RESPONSE_DATA_SLEEP_TIME)
		}

		// 开始解析
		t := &pool.Task{}
		t.TaskFunc = doparse
		t.Param = []interface{}{parseFunc, resp, p}
		err := g.Pool.Submit(t)
		if err != nil {
			logs.Logger.Error("g.Pool.Submit", zap.Error(err))
		}
	}

	// todo 满足某些条件关闭 group 和 spider
	//g.Status = false
}

// 存储新数据
func storage(g *pool.Group, p pipeline.Pipeline) {
	for {
		if p.Empty() {
			// todo 累计超过某一值
			time.Sleep(NO_PIPELINE_DATA_SLEEP_TIME)
			continue
		}
		t := &pool.Task{}
		t.TaskFunc = func(i []interface{}) {
			err := p.Save()
			if err != nil {
				logs.Logger.Error("storage.Save", zap.Error(err))
			}
			return
		}
		t.Param = []interface{}{}
		err := g.Pool.Submit(t)
		if err != nil {
			logs.Logger.Error("g.Pool.Submit", zap.Error(err))
		}
	}

	// todo 满足某些条件关闭 group 和 spider
	//g.Status = false
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
	p := params[2].(pipeline.Pipeline)

	// 处理得出的结果
	result := f(r)

	// 传输到通道保存
	p.Push(result)
}
