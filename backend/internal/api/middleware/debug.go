package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// RequestResponseLogger logs detailed request and response information for debugging
func RequestResponseLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		var statusColor, methodColor, resetColor string
		if param.IsOutputColor() {
			statusColor = param.StatusCodeColor()
			methodColor = param.MethodColor()
			resetColor = param.ResetColor()
		}

		if param.Latency > time.Minute {
			param.Latency = param.Latency.Truncate(time.Second)
		}

		return fmt.Sprintf("[GIN] %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			statusColor, param.StatusCode, resetColor,
			param.Latency,
			param.ClientIP,
			methodColor, param.Method, resetColor,
			param.Path,
			param.ErrorMessage,
		)
	})
}

// DebugRequestResponse logs full request and response details in debug mode
func DebugRequestResponse() gin.HandlerFunc {
	return func(c *gin.Context) {
		if logrus.GetLevel() != logrus.DebugLevel {
			c.Next()
			return
		}

		// Log request
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		logrus.WithFields(logrus.Fields{
			"method":  c.Request.Method,
			"path":    c.Request.URL.Path,
			"query":   c.Request.URL.RawQuery,
			"headers": c.Request.Header,
			"body":    string(requestBody),
			"ip":      c.ClientIP(),
		}).Debug("Incoming request")

		// Capture response
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		start := time.Now()
		c.Next()
		latency := time.Since(start)

		// Log response
		responseBody := blw.body.String()
		
		logEntry := logrus.WithFields(logrus.Fields{
			"status":   c.Writer.Status(),
			"latency":  latency,
			"size":     c.Writer.Size(),
			"headers":  c.Writer.Header(),
		})

		// Try to format JSON response for better readability
		if json.Valid([]byte(responseBody)) {
			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, []byte(responseBody), "", "  "); err == nil {
				logEntry = logEntry.WithField("body", prettyJSON.String())
			} else {
				logEntry = logEntry.WithField("body", responseBody)
			}
		} else {
			logEntry = logEntry.WithField("body", responseBody)
		}

		logEntry.Debug("Outgoing response")
	}
}
