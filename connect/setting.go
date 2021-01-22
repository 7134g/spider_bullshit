package connect

import "time"

var (
	// 调度器取不到 resp 时（即队列为空），且值为nil，若中间件不处理，沉睡时间
	NO_RESPONSE_DATA_SLEEP_TIME = time.Millisecond * 500
	//
)
