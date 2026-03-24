// Package middleware 限流中间件
package middleware

import (
	"github.com/MJ-9527/GoMind/config"
	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth_gin"
	"github.com/gin-gonic/gin"
)

func RateLimitMiddleware() gin.HandlerFunc {
	maxRequests := config.GlobalConfig.RateLimit.MaxRequests
	limiter := tollbooth.NewLimiter(float64(maxRequests), nil)
	limiter.SetIPLookups([]string{"RemoteAddr", "X-Forwarded-For"})
	return tollbooth_gin.LimitHandler(limiter)
}
