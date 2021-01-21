package test

import (
	"bullshit"
)

func main() {
	// 设置全局默认
	p := bullshit.NewCtrlParams()

	// 总控制器
	ctrl := bullshit.NewController(p)
	// 通过基础设置生成爬虫
	s1 := ctrl.NewSpider()
	/*
		设置该爬虫参数
		s1.ProxyValue = "11111111"
		...
	*/
	// 装载爬虫
	ctrl.LoadSpider(s1)
	// 启动调度器，开始执行爬虫
	ctrl.Run()
}
