package logger

import (
	"context"
	"io"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type handlerLogConfig struct {
	logLevel logrus.Level
}

type HandlerLogOption func(conf *handlerLogConfig)

func WithGinDebug(level logrus.Level) HandlerLogOption {
	return func(conf *handlerLogConfig) {
		conf.logLevel = level
	}
}

type handlerLogger struct {
	*logrus.Logger
}

var _ io.Writer = (*handlerLogger)(nil)

var logPool *sync.Pool

func init() {
	logPool = &sync.Pool{
		New: func() interface{} {
			log := logrus.New()
			log.SetFormatter(&logrus.TextFormatter{
				TimestampFormat: TimestampFormat,
			})
			return &handlerLogger{
				Logger: log,
			}
		},
	}
}

func NewHandlerLogger() *handlerLogger {
	return logPool.Get().(*handlerLogger)
}

func (l *handlerLogger) Write(p []byte) (n int, err error) {
	p = []byte(strings.ReplaceAll(string(p), "[GIN-debug] ", "[API-RBAC-debug] "))
	l.Logger.Debug(string(p))
	return len(p), nil
}

func RestLogger(opts ...HandlerLogOption) gin.HandlerFunc {
	conf := &handlerLogConfig{
		logLevel: logrus.InfoLevel,
	}
	for _, opt := range opts {
		opt(conf)
	}
	return func(c *gin.Context) {
		log := NewHandlerLogger()
		log.SetLevel(conf.logLevel)

		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		ctx := c.Request.Context()
		mpHeader := c.Request.Header.Clone()
		newCtx := ToContext(ctx, log)
		*c.Request = *c.Request.WithContext(newCtx)
		c.Next()

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
		fields := logrus.Fields{
			"key":      "RBNS",
			"status":   param.StatusCode,
			"latency":  param.Latency,
			"clientIP": param.ClientIP,
			"method":   param.Method,
			"path":     param.Path,
			"Ua":       param.Request.UserAgent(),
		}
		if conf.logLevel == logrus.DebugLevel {
			for key, value := range mpHeader {
				if len(value) >= 0 {
					log.Debugf("req_%s: %v", key, value)
				}
			}
		}
		if conf.logLevel == logrus.DebugLevel {
			mpHeader := c.Writer.Header().Clone()
			for key, value := range mpHeader {
				if len(value) >= 0 {
					log.Debugf("res_%s: %v", key, value)
				}
			}
		}
		log.WithFields(fields).Info("incoming request")
		logPool.Put(log)
	}
}

func codeFunc(err error) codes.Code {
	return status.Code(err)
}

func GrpcLogger(opts ...HandlerLogOption) grpc.UnaryServerInterceptor {
	conf := &handlerLogConfig{
		logLevel: logrus.InfoLevel,
	}
	for _, opt := range opts {
		opt(conf)
	}
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		log := NewHandlerLogger()
		log.SetLevel(conf.logLevel)
		fullMethodString := info.FullMethod
		service := path.Dir(fullMethodString)[1:]
		method := path.Base(fullMethodString)
		start := time.Now()
		fields := logrus.Fields{
			"grpc.key":     "[API-RBAC]",
			"grpc.start":   start.Format(TimestampFormat),
			"grpc.service": service,
			"grpc.method":  method,
		}
		if d, ok := ctx.Deadline(); ok {
			fields["grpc.request.deadline"] = d.Format(TimestampFormat)
		}
		newCtx := ToContext(ctx, log)
		resp, err = handler(newCtx, req)
		code := codeFunc(err)
		timestamp := time.Now()
		latency := timestamp.Sub(start)
		fields["grpc.code"] = code.String()
		fields["grpc.latency"] = latency
		if err != nil {
			fields[logrus.ErrorKey] = err
		}
		log.WithContext(newCtx).WithFields(fields).Info("finished unary call with code " + code.String())
		logPool.Put(log)
		return
	}
}

type ctxLoggerMarker struct{}

var logKey = &ctxLoggerMarker{}

func FromContext(ctx context.Context) *handlerLogger {
	if val, ok := ctx.Value(logKey).(*handlerLogger); ok && val != nil {
		return val
	}
	log := NewHandlerLogger()
	log.SetLevel(logrus.DebugLevel)
	return log
}

func ToContext(ctx context.Context, log *handlerLogger) context.Context {
	return context.WithValue(ctx, logKey, log)
}
