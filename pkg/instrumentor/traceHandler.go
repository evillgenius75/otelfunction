package instrumentor

import (
	"context"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

type Flush interface {
	ForceFlush(context.Context) error
}

type HttpHandler = func(w http.ResponseWriter, r *http.Request)

func InstrumentedHandler(functionName string, function HttpHandler, flusher Flush) HttpHandler {
	opts := []trace.SpanStartOption{
		// customizable span attributes
		trace.WithAttributes(semconv.FaaSTriggerHTTP),
	}

	// create instrumented handler
	handler := otelhttp.NewHandler(
		http.HandlerFunc(function), functionName, otelhttp.WithSpanOptions(opts...),
	)

	return func(w http.ResponseWriter, r *http.Request) {
		// call the actual handler
		handler.ServeHTTP(w, r)

		// flush spans
		flusher.ForceFlush(r.Context())
	}
}
