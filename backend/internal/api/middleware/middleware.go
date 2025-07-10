package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

// RateLimiter represents a rate limiter middleware
type RateLimiter struct {
	redis     *redis.Client
	rateLimit int
	window    time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(redisClient *redis.Client, rateLimit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		redis:     redisClient,
		rateLimit: rateLimit,
		window:    window,
	}
}

// Middleware returns the rate limiting middleware
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		key := fmt.Sprintf("rate_limit:%s", clientIP)
		
		ctx := context.Background()
		
		// Get current count
		current, err := rl.redis.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			logrus.Errorf("Rate limiter Redis error: %v", err)
			c.Next()
			return
		}
		
		// Check if limit exceeded
		if current >= rl.rateLimit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
				"retry_after": rl.window.Seconds(),
			})
			c.Abort()
			return
		}
		
		// Increment counter
		pipe := rl.redis.Pipeline()
		pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, rl.window)
		
		if _, err := pipe.Exec(ctx); err != nil {
			logrus.Errorf("Rate limiter Redis pipeline error: %v", err)
		}
		
		c.Next()
	}
}

// Logger returns a gin.LoggerWithFormatter middleware with custom format
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logrus.WithFields(logrus.Fields{
			"status":      param.StatusCode,
			"method":      param.Method,
			"path":        param.Path,
			"ip":          param.ClientIP,
			"user_agent":  param.Request.UserAgent(),
			"latency":     param.Latency,
			"time":        param.TimeStamp.Format(time.RFC3339),
		}).Info("HTTP Request")
		
		return ""
	})
}

// Recovery returns a gin.Recovery middleware
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logrus.WithFields(logrus.Fields{
			"panic": recovered,
			"path":  c.Request.URL.Path,
			"ip":    c.ClientIP(),
		}).Error("Panic recovered")
		
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
	})
}

// Security middleware adds various security headers
func Security() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Next()
	}
}

// RequestID middleware adds a unique request ID to each request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		
		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}

// generateRequestID generates a simple request ID
func generateRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
