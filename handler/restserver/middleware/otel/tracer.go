package otel

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/n-creativesystem/rbns/handler/restserver/contexts"
	"github.com/n-creativesystem/rbns/logger"
	"github.com/n-creativesystem/rbns/tracer"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	tracerKey = "otel-middleware"
)

func Middleware(serviceName string) contexts.HandlerFunc {
	return func(c *contexts.Context) error {
		savedCtx := c.Request.Context()
		defer func() {
			c.Request = c.Request.WithContext(savedCtx)
		}()
		ctx := otel.GetTextMapPropagator().Extract(savedCtx, propagation.HeaderCarrier(c.Request.Header))
		opts := []trace.SpanStartOption{
			trace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", c.Request)...),
			trace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(c.Request)...),
			trace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(serviceName, c.FullPath(), c.Request)...),
			trace.WithSpanKind(trace.SpanKindServer),
		}
		spanName := c.FullPath()
		if spanName == "" {
			spanName = fmt.Sprintf("HTTP %s route not found", c.Request.Method)
		}
		ctx, span := tracer.Start(ctx, spanName, opts...)
		defer span.End()

		c.Request = c.Request.WithContext(ctx)

		c.Context.Next()

		status := c.Writer.Status()
		attrs := semconv.HTTPAttributesFromHTTPStatusCode(status)
		spanStatus, spanMessage := semconv.SpanStatusFromHTTPStatusCode(status)
		span.SetAttributes(attrs...)
		span.SetStatus(spanStatus, spanMessage)
		if len(c.Errors) > 0 {
			span.SetAttributes(attribute.String("gin.errors", c.Errors.String()))
		}
		return nil
	}
}

func RestLogger(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		span := trace.SpanFromContext(c.Request.Context())
		defer span.End()
		logger.SetLogger(c, log)
		w := &writer{
			ResponseWriter: c.Writer,
			buffer:         []byte{},
		}
		c.Writer = w
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		fields := make([]interface{}, 0, 100)
		mpHeader := c.Request.Header.Clone()
		for key, value := range mpHeader {
			if len(value) >= 0 {
				k := fmt.Sprintf("req-%s", strings.ToLower(key))
				v := strings.ToLower(strings.Join(value, ", "))
				fields = append(fields, []interface{}{k, v}...)
				span.SetAttributes(attribute.StringSlice(k, value))
			}
		}
		c.Next()
		for i, err := range c.Errors {
			log.ErrorWithContext(c, err, fmt.Sprintf("idx: %d error: %v", i, err))
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
		for k, v := range mp {
			switch value := v.(type) {
			case string:
				span.SetAttributes(attribute.String(k, value))
			case int:
				span.SetAttributes(attribute.Int(k, value))
			case time.Duration:
				span.SetAttributes(attribute.String(k, value.String()))
			}
			fields = append(fields, []interface{}{k, v}...)
		}
		mpHeader = c.Writer.Header().Clone()
		for key, value := range mpHeader {
			if len(value) >= 0 {
				k := fmt.Sprintf("res-%s", strings.ToLower(key))
				v := strings.ToLower(strings.Join(value, ", "))
				fields = append(fields, []interface{}{k, v}...)
			}
		}
		if c.Writer.Status() > 299 {
			log.Info(w.String(), fields...)
		} else {
			log.Info("incoming request", fields...)
		}
	}
}

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
