package tracing

import (
	"net/http"
)

type TracingTransport struct {
	Transport http.RoundTripper
}

func NewTracingTransport(transport http.RoundTripper) *TracingTransport {
	return &TracingTransport{
		Transport: transport,
	}
}

func (t *TracingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	span := FromContext(ctx)
	if span != nil {
		req.Header.Set("X-Trace-ID", span.TraceID)
		req.Header.Set("X-Span-ID", span.SpanID)
	}
	return t.Transport.RoundTrip(req)
}
