package pipeline

import (
	"container/list"
	"gorm.io/gorm"
	"sync"
)

type Pipeline interface {
	Push(v interface{})
	Save() error
	Empty() bool
	Len() int
}

type PipelineBase struct {
	list *list.List
	mux  sync.Mutex
}

// 获取长度
func (p *PipelineBase) Len() int {
	length := p.list.Len()
	return length
}

func (p *PipelineBase) Empty() bool {
	p.mux.Lock()
	defer p.mux.Unlock()

	if p.Len() == 0 {
		return true
	}
	return false
}

type SQLitePipeline struct {
	PipelineBase
	db *gorm.DB
}

func (p *SQLitePipeline) Push(v interface{}) {
	p.mux.Lock()
	p.list.PushFront(v)
	p.mux.Unlock()
}

func (p *SQLitePipeline) Save() error {
	p.mux.Lock()
	defer p.mux.Unlock()
	var inter interface{}

	if p.Len() > 0 {
		inter = p.list.Remove(p.list.Back())
		err := p.db.Create(inter).Error
		if err != nil {
			return err
		}
	}

	return nil
}
