package pool

import (
	"context"
	"time"
)

type Group struct {
	Name        string             // 组名
	workerCount int                // 工人数
	Pool        *Pool              // 工作组
	StopFunc    context.CancelFunc // 关闭工作组信号
}

func AddNewGroup(name string, workerCount int) *Group {
	p, cencel := NewPool(int32(workerCount), false, 30*time.Second)
	g := &Group{
		Name:        name,
		workerCount: workerCount,
		Pool:        p,
		StopFunc:    cencel,
	}
	return g
}
