package connect

import "time"

var (
	// 调度器取不到 requests 时（即队列为空），且值为nil，若中间件不处理，沉睡时间
	NO_REQUESTS_DATA_SLEEP_TIME = time.Millisecond * 500
	// 调度器取不到 response 时（即队列为空），且值为nil，若中间件不处理，沉睡时间
	NO_RESPONSE_DATA_SLEEP_TIME = time.Millisecond * 500
	// 通道数据为空
	NO_PIPELINE_DATA_SLEEP_TIME = time.Millisecond * 500
)
