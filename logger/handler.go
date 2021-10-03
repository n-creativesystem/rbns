package logger

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type writer struct {
	gin.ResponseWriter
	buffer []byte
}

func (w *writer) String() string {
	return string(w.buffer)
}

func (w *writer) Write(p []byte) (int, error) {
	n, err := w.ResponseWriter.Write(p)
	if err != nil {
		return n, err
	}
	w.buffer = append(w.buffer, p...)
	return n, err
}

func RestLogger(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		SetLogger(c, logger)
		w := &writer{
			ResponseWriter: c.Writer,
			buffer:         []byte{},
		}
		c.Writer = w
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		fields := logrus.Fields{}
		mpHeader := c.Request.Header.Clone()
		for key, value := range mpHeader {
			if len(value) >= 0 {
				k := fmt.Sprintf("req-%s", strings.ToLower(key))
				v := strings.ToLower(strings.Join(value, ", "))
				fields[k] = v
			}
		}
		c.Next()
		for i, err := range c.Errors {
			logger.Errorf("idx: %d error: %v", i, err)
		}
		param := gin.LogFormatterParams{
			Request: c.Request,
			Keys:    c.Keys,
		}

		param.TimeStamp = time.Now()
		param.Latency = param.TimeStamp.Sub(start)

		param.ClientIP = c.ClientIP()
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()

		param.BodySize = c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		param.Path = path
		mp := map[string]interface{}{
			"key":      "RBNS",
			"status":   param.StatusCode,
			"latency":  param.Latency,
			"clientIP": param.ClientIP,
			"method":   param.Method,
			"path":     param.Path,
			"Ua":       param.Request.UserAgent(),
		}
		for key, value := range mp {
			fields[key] = value
		}
		mpHeader = c.Writer.Header().Clone()
		for key, value := range mpHeader {
			if len(value) >= 0 {
				k := fmt.Sprintf("res-%s", strings.ToLower(key))
				v := strings.ToLower(strings.Join(value, ", "))
				fields[k] = v
			}
		}
		if c.Writer.Status() > 299 {
			logger.WithFields(fields).Info(w.String())
		} else {
			logger.WithFields(fields).Info("incoming request")
		}
	}
}

type writeLogger struct {
	log *logrus.Logger
}

func (log *writeLogger) Write(p []byte) (int, error) {
	log.log.Info(string(p))
	return len(p), nil
}

func NewWriter(log *logrus.Logger) io.Writer {
	return &writeLogger{
		log: log,
	}
}

const (
	logKey = "logger"
)

func GetLogger(c *gin.Context) *logrus.Logger {
	if v, ok := c.Get(logKey); ok {
		return v.(*logrus.Logger)
	}
	return logrus.StandardLogger()
}

func SetLogger(c *gin.Context, logger *logrus.Logger) {
	c.Set(logKey, logger)
}
