package tracing

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type traceContextKey struct{}

type Span struct {
	TraceID    string
	SpanID     string
	ParentSpan *Span
}

type TraceInit struct {
	ServiceName string
	Logpath     string
}

func NewSpan(traceID string, parent *Span) *Span {
	return &Span{
		TraceID:    traceID,
		SpanID:     uuid.New().String(),
		ParentSpan: parent,
	}
}

func FromContext(ctx context.Context) *Span {
	span, _ := ctx.Value(traceContextKey{}).(*Span)
	return span
}

func NewContext(ctx context.Context, span *Span) context.Context {
	return context.WithValue(ctx, traceContextKey{}, span)
}

func TracingMiddleware(init TraceInit) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := os.OpenFile(init.Logpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		log.SetOutput(file)
		log.SetFlags(log.Ldate | log.Ltime | log.LUTC)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}

		traceID := c.Request.Header.Get("X-Trace-ID")
		spanID := c.Request.Header.Get("X-Span-ID")

		var span *Span
		var parentSpan *Span

		if traceID != "" && spanID != "" {
			parentSpan = &Span{TraceID: traceID, SpanID: spanID}
			span = NewSpan(traceID, parentSpan)
		} else {
			span = NewSpan(uuid.New().String(), nil)
			c.Request.Header.Set("X-Trace-ID", span.TraceID)
			c.Request.Header.Set("X-Span-ID", span.SpanID)
		}

		ctx := NewContext(c.Request.Context(), span)
		c.Request = c.Request.WithContext(ctx)

		start := time.Now()
		c.Next()
		duration := time.Since(start)

		log.Printf("Service: %s, TraceID: %s, SpanID: %s, ParentSpanID: %s, Method: %s, Path: %s, Duration: %d, Status: %d",
			init.ServiceName, span.TraceID, span.SpanID, parentSpanID(span), c.Request.Method, c.Request.URL.Path, duration.Milliseconds(), c.Writer.Status())

		c.Writer.Header().Set("X-Trace-ID", span.TraceID)
		c.Writer.Header().Set("X-Span-ID", span.SpanID)
	}
}

func parentSpanID(span *Span) string {
	if span != nil && span.ParentSpan != nil {
		return span.ParentSpan.SpanID
	}
	return ""
}
