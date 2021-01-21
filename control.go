package bullshit

import (
	"bullshit/common/logs"
	"bullshit/common/serializer"
	"time"
)

var (
	controller *Controller
)

// 全局配置
type CtrlParams struct {
	BaseMiddlewares []SpiderOption
	ConcurrentMax   uint          // 最大并发数
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
}

func NewCtrlParams() *CtrlParams {
	return &CtrlParams{
		Logger: logs.NewLogger(),
	}
}

func (p *CtrlParams) LoadMiddleware(sp ...SpiderOption) {
	p.BaseMiddlewares = append(p.BaseMiddlewares, sp...)
}

// 控制器
type Controller struct {
	Spiders []*Spider
	Sche    *Scheduler
	Params  *CtrlParams
}

func NewController(p *CtrlParams) *Controller {
	sche := NewScheduler()
	c := &Controller{
		Sche: sche,
	}
	for _, f := range c.Params.BaseMiddlewares {
		f(c)
	}
	c.Sche.ConcurrentMax = p.ConcurrentMax
	controller = controller
	return c
}

func (c *Controller) NewSpider() *Spider {

	s := &Spider{}
	serializer.StructValue(c, s)
	return s
}

// 加载爬虫
func (c *Controller) LoadSpider(cs *Spider) {
	for _, s := range c.Spiders {
		if s.Name == cs.Name {
			return
		}
	}
	c.Spiders = append(c.Spiders, cs)
}

func IgnoreRobotsTxt() func(*Spider) {
	return func(c *Spider) {
		c.IgnoreRobotsTxt = true
	}
}

func LoggerSysoutput() func(*Spider) {
	return func(c *Spider) {
		l := &logs.Logger{
			Silent:     &c.SystemOutput,
			Level:      &c.LogLevel,
			Path:       &c.LogPath,
			FileStatus: &c.LogFile,
		}
		l.Init()
		c.Logger = l
	}
}

func (c *Controller) Run() {

}
