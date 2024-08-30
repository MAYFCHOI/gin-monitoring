package tracing

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type traceContextKey struct{}

type Span struct {
	TraceID string
	SpanID  string
}

func NewSpan() *Span {
	return &Span{
		TraceID: uuid.New().String(),
		SpanID:  uuid.New().String(),
	}
}

func FromContext(ctx context.Context) *Span {
	span, _ := ctx.Value(traceContextKey{}).(*Span)
	return span
}

func NewContext(ctx context.Context, span *Span) context.Context {
	return context.WithValue(ctx, traceContextKey{}, span)
}

func TracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.Request.Header.Get("X-Trace-ID")
		spanID := c.Request.Header.Get("X-Span-ID")

		var span *Span
		if traceID != "" && spanID != "" {
			span = &Span{TraceID: traceID, SpanID: spanID}
		} else {
			span = NewSpan()
		}

		ctx := NewContext(c.Request.Context(), span)
		c.Request = c.Request.WithContext(ctx)

		// 요청에 Trace ID와 Span ID 헤더 추가
		c.Request.Header.Set("X-Trace-ID", span.TraceID)
		c.Request.Header.Set("X-Span-ID", span.SpanID)

		start := time.Now()
		c.Next()
		duration := time.Since(start)

		log.Printf("TraceID: %s, SpanID: %s, Method: %s, Path: %s, Duration: %s, Status: %d",
			span.TraceID, span.SpanID, c.Request.Method, c.Request.URL.Path, duration, c.Writer.Status())
	}
}
