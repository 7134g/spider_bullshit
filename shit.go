package bullshit

import (
	"bullshit/common/logs"
	"bullshit/connect"
	"sync"
	"time"
)

//	SpiderName string // 爬虫名字
//	//Spider            *colly.Spider                      // 爬虫控制器
//	//Pool                 *pool.Queue                           // 任务队列
//	//WG                   *sync.WaitGroup // 任务锁
//	ConcurrentCount int  // 并发数
//	Duplicate       bool // 去重
//	//continuousRepetition int32           // 连续重复任务数
//	//Proxy                func(*http.Request) (*url.URL, error) // 开启代理
//	ProxyValue   string        // 代理ip值
//	HttpTimeout  time.Duration // 连接超时时间
//	SpiderStatus bool          // 爬虫状态
//	WaitTime     time.Duration // 等待时间
//
//	// 不能配置字段
//	synchronous         sync.Mutex // 爬虫锁
//	statusLock          sync.Mutex // 爬虫开关锁
//
//
//	// middleware
//	errorCount          int32      // 连续错误次数
//	retryCount          int32      // 重试上限
//	spiderSleepMaxCount int        // 爬虫最长等待时间 默认为 50秒，检查100次, 间隔 500ms

// base
type Spider struct {
	Name string // 爬虫名字

	ConcurrentCount int           // 并发数，chan ?
	Duplicate       bool          // 去重
	ProxyValue      string        // 代理ip值
	Timeout         time.Duration // 连接超时时间
	WaitTime        time.Duration // 请求间隔时间
	SpiderStatus    bool          // 爬虫状态,关闭,打开
	OnCookieJar     bool          // 是否开启cookiesJar管理cookie
	DefaultHeader   bool          // 默认头
	IgnoreRobotsTxt bool          // Robots.txt

	Logger       *logs.Logger // 日志对象
	LogLevel     logs.Level   // 日志等级
	LogPath      string       // log日志位置
	LogFile      bool         // 是否开启日志文件
	SystemOutput bool         // 禁止标准输出

	statusLock             *sync.RWMutex // 爬虫开关锁
	requestCallbacks       []RequestCallback
	responseCallbacks      []ResponseCallback
	baseMiddlewareCallback []BaseMiddlewareCallback
}

type SpiderOption func(*Spider)

// RequestsMiddlerware is a type alias for OnRequest callback functions
type RequestCallback func(*connect.Requests)

// ResponseMiddlerware is a type alias for OnResponse callback functions
type ResponseCallback func(*connect.Response)

type BaseMiddlewareCallback func(*connect.Response)

// 构建请求前
func (c *Spider) BeforeRequests() {

}

// 对请求操作
func (c *Spider) OnRequests() {

}

// 构建响应前
func (c *Spider) BeforeResponse() {

}

// 对响应操作
func (c *Spider) OnResponse() {

}

// 设置超时
func (c *Spider) SetTimeout(t time.Duration) {
	c.Timeout = t
}

// 设置代理
func (c *Spider) SetProxy(ip string) {
	c.ProxyValue = ip
}

// 设置等待时间，此时为单线程
func (c *Spider) SetWaitTime(t time.Duration) {
	c.WaitTime = t
	c.ConcurrentCount = 1
}

// 停止收集器
func (c *Spider) Stop() {
	c.statusLock.Lock()
	defer c.statusLock.Unlock()
	c.SpiderStatus = false
}
