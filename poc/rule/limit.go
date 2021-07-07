package rule

import (
	"context"
	"github.com/jweny/pocassist/pkg/conf"
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

// 限制速率
var limiter *rate.Limiter

// 批量资产测试并发channal
var OriginalReqChannel chan *http.Request

func InitOreqChannel()  {
	concurrent := 10
	if conf.GlobalConfig.PluginsConfig.Concurrent != 0 {
		concurrent = conf.GlobalConfig.PluginsConfig.Concurrent
	}
	OriginalReqChannel = make(chan *http.Request, concurrent)
}

func InitRate() {
	maxQps := conf.GlobalConfig.HttpConfig.MaxQps
	parallel := conf.GlobalConfig.PluginsConfig.Parallel

	limit := rate.Every(time.Duration(maxQps) * time.Millisecond)
	// 第二个参数 和 并发加载的 plugin 数匹配
	limiter = rate.NewLimiter(limit, parallel)
}

func LimitWait() {
	limiter.Wait(context.Background())
}

