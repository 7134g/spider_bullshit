package bullshit

import (
	"bullshit/pool"
	"context"
	"time"
)

// 调度器
type Scheduler struct {
	Collectors    []*Spider // 所有收集器
	ConcurrentMax uint      // 最大并发数
	Pool          *pool.Pool
	CencelPool    context.CancelFunc
}

func NewScheduler() *Scheduler {
	p, cencel := pool.NewPool(1, true, time.Second*1)
	return &Scheduler{
		Collectors: []*Spider{},
		Pool:       p,
		CencelPool: cencel,
	}
}

func (s *Scheduler) Run() {

}
