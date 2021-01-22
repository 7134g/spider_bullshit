package connect

import "testing"

type BiliSpider struct {
	Spider
}

func TestScheduler_Excute(t *testing.T) {
	sche := NewScheduler()
	s1 := NewSpider("s1")
	sche.AddSpider(s1)

	s2 := &BiliSpider{}
	s2.RequQueue = NewRequCh(10000)
	s2.ParseFunc = func(resp *Response) {
	}

	sche.Excute()
}
