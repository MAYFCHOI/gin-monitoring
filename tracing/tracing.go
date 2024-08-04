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

// 분산 추적 미들웨어
func TracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		span := NewSpan()
		ctx := NewContext(c.Request.Context(), span)
		c.Request = c.Request.WithContext(ctx)

		start := time.Now()
		c.Next()
		duration := time.Since(start)
		log.Printf("TraceID: %s, SpanID: %s, Method: %s, Path: %s, Duration: %s, Status: %d",
			span.TraceID, span.SpanID, c.Request.Method, c.Request.URL.Path, duration, c.Writer.Status())
	}
}
