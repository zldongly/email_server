package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

type Logger interface {
	Info(msg ...interface{})
}

func Log(log Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		f := format{
			ClientIP: c.ClientIP(),
			Method : c.Request.Method,
			Path : c.Request.URL.Path,
		}
		start := time.Now()
		raw := c.Request.URL.RawQuery
		if raw != "" {
			f.Path += "?" + raw
		}

		c.Next()

		f.StatusCode = c.Writer.Status()
		f.BodySize = c.Writer.Size()
		f.Latency = time.Now().Sub(start) // 处理时间

		log.Info(f.String())
	}
}

type format struct {
	ClientIP   string
	Path       string
	Method     string
	BodySize      int
	Latency    time.Duration
	StatusCode int
}

func (f *format) String() string {
	return fmt.Sprintf("client_ip = %s, path = %s, method = %s, body_size = %d, latency = %v, status_code = %d",
		f.ClientIP, f.Path, f.Method, f.BodySize, f.Latency, f.StatusCode)
}
