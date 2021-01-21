#### base middlerware
1. duplicate local
2. duplicate remote
3. cookies
4. ua
5. error
6. log

#### base pipeline
1. local file


#### default setting
	ConcurrentCount int  // 并发数
	Duplicate       bool // 去重
	ProxyValue   string  // 代理ip值
	Timeout  time.Duration // 连接超时时间
	WaitTime     time.Duration// 请求间隔时间
	SpiderStatus bool          // 爬虫状态,关闭,打开





#### Life Cycle
1. use -> create requests
2. Mrequests deal requests
3. Crequests push to Scheduler
4. Scheduler duplicate removal Crequests
5. Scheduler will Crequests store to local file or queue
6. get queue Crequests to take and initiate HTTP
7. waiting for Crequests done
8. Scheduler will push response to Mresponse
9. Mresponse push Cresponse to Scheduler
10. Scheduler Cresponse -> use
11. use deal Cresponse -> Scheduler
12. Scheduler push to pipeline


Crequests: complete requests  
Cresponse: complete response  
Mrequests: middlerware requests  
Mresponse: middlerware response  

#### error
Once middlerware intercepts an error, tell the scheduler to stop processing the request.

#### init
	controller
	Scheduler
	Mrequests
		|-----M1
		|-----M2
		|	...
	Mresponse
		|-----M1
		|-----M2
		|	...
	pipeline
