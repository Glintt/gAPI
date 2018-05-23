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

	currentRequestMetricName := GetRateLimitingMetricName(c)

	CreateRateLimitingIfNotExists(currentRequestMetricName)

	rateStatus := RateLimitingStatusForRequest(currentRequestMetricName)

	if IsRateLimitExceeded(rateStatus) {
		http.Response(c, `{"error":true, "msg": "Rate limiting exceeded."}`, 429, c.Request.URI().String())
		c.Abort()
		return nil
	}
	
	limits[currentRequestMetricName] = RateStatus{NumberRequests: (rateStatus.NumberRequests + 1), ExpirationTime: rateStatus.ExpirationTime}
	return nil
}

func CreateRateLimitingIfNotExists(currentRequestMetricName string) {
	if _, ok := limits[currentRequestMetricName]; ok == false {
		expirationTime := RateLimitingExpirationTime()
		limits[currentRequestMetricName] = RateStatus{NumberRequests:1, ExpirationTime: expirationTime}
	}
}

func RateLimitingStatusForRequest(currentRequestMetricName string) RateStatus {
	currentNumberRequests := limits[currentRequestMetricName].NumberRequests
	currentExpirationTime := limits[currentRequestMetricName].ExpirationTime

	if currentExpirationTime < utils.CurrentTimeMilliseconds() {
		currentExpirationTime = RateLimitingExpirationTime()
		currentNumberRequests = 0
	}

	return RateStatus{NumberRequests: currentNumberRequests, ExpirationTime: currentExpirationTime}
}

func IsRateLimitExceeded(rateStatus RateStatus) bool {
	if rateStatus.ExpirationTime > utils.CurrentTimeMilliseconds() && rateStatus.NumberRequests < limiter.Limit {
		return false
	}
	return true
}

func GetRateLimitingMetricName(c *routing.Context) string {
	return c.RemoteAddr().String()
}