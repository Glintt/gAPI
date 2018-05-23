package ratelimiting

import (
	"gAPIManagement/api/config"
	"gAPIManagement/api/http"
	"gAPIManagement/api/utils"
	"time"
	
	"github.com/qiangxue/fasthttp-routing"
)

type LimiterRate struct {
	Period time.Duration
	Limit int
}

type RateStatus struct {
	NumberRequests int
	ExpirationTime int64
}


var limiter config.GApiRateLimitingConfig 
var limits map[string]RateStatus

func InitRateLimiting() {
	limiter = config.GApiConfiguration.RateLimiting
	limits = make(map[string]RateStatus)
}

func RateLimitingExpirationTime() int64 {
	return utils.CurrentTimeMilliseconds() + (limiter.Period * 60 * 1000)
}

func RateLimiting(c *routing.Context) error {
	if ! limiter.Active {
		return nil
	}
	
	if _, ok := limits[c.RemoteAddr().String()]; ok == false {
		expirationTime := RateLimitingExpirationTime()
		limits[c.RemoteAddr().String()] = RateStatus{NumberRequests:1, ExpirationTime: expirationTime}
		return nil
	}

	currentNumberRequests := limits[c.RemoteAddr().String()].NumberRequests
	currentExpirationTime := limits[c.RemoteAddr().String()].ExpirationTime

	if currentExpirationTime < utils.CurrentTimeMilliseconds() {
		currentExpirationTime = RateLimitingExpirationTime()
		currentNumberRequests = 0
	}

	if currentExpirationTime > utils.CurrentTimeMilliseconds() && currentNumberRequests < limiter.Limit {
		currentNumberRequests = currentNumberRequests + 1
		limits[c.RemoteAddr().String()] = RateStatus{NumberRequests: currentNumberRequests, ExpirationTime: currentExpirationTime}
		return nil
	}

	http.Response(c, `{"error":true, "msg": "Rate limiting exceeded."}`, 429, c.Request.URI().String())
	c.Abort()
	return nil
}