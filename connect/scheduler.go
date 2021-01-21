package connect

import (
	"bullshit/common/logs"
	"bullshit/pool"
	"go.uber.org/zap"
	"net/http"
	"sync"
)

type Scheduler struct {
	lock       sync.RWMutex
	SpiderList []SpiderBase  // 所有爬虫类型
	Groups     []*pool.Group // 所有工作组

	status                  status // 调度器状态
	baseRequestsMiddlerware []RequestsMiddlerware
	baseResponseMiddlerware []ResponseMiddlerware
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		SpiderList: make([]SpiderBase, 0),
		Groups:     make([]*pool.Group, 0),
	}
}

func (s *Scheduler) AddSpider(spider SpiderBase) {
	s.lock.Lock()
	defer s.lock.Unlock()
	// 装载共有模块
	spider.allBaseMiddlerware(s.baseRequestsMiddlerware, s.baseResponseMiddlerware)
	s.SpiderList = append(s.SpiderList, spider)
}

func (s *Scheduler) Excute() {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.status == running {
		return
	}
	s.status = running

	// 为每一个爬虫创建工作组, 添加任务执行到工作组
	for _, spider := range s.SpiderList {
		name, count := spider.groupBase()
		group := pool.AddNewGroup(name, count)
		s.Groups = append(s.Groups, group)
		go monitor(group, spider)
	}

}

// 监控
func monitor(g *pool.Group, s SpiderBase) {
	for {
		select {
		case req := <-s.pushRequestToSche():
			t := &pool.Task{}
			t.TaskFunc = dohttp
			client, respResult := s.spiderBase()
			t.Param = []interface{}{req, client, respResult}
			err := g.Pool.Submit(t)
			if err != nil {
				logs.Logger.Error("g.Pool.Submit", zap.Error(err))
			}
		}
	}

}

func dohttp(params []interface{}) {
	req := params[0].(*Requests)
	client := params[1].(*http.Client)
	respResult := params[2].(*[]Response)

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

	response, err := NewResponse(req.Name, resp)
	if err != nil {
		logs.Logger.Error("read body error", zap.Error(err))
		return
	}

	*respResult = append(*respResult, *response)

}
