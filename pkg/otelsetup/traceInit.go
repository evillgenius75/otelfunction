package otelsetup

import (
	"context"
	"log"
	"os"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/otel"

	//"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	//"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

func InitTracing(serviceName string, serviceVersion string) *sdk.TracerProvider {
	//client := otlptracehttp.NewClient()
	ctx := context.Background()
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	exporter, err := texporter.New(texporter.WithProjectID(projectID))
	if err != nil {
		log.Fatalf("texporter.New: %v", err)
	}

	// create resource
	r, _ := resource.New(ctx,
		resource.WithTelemetrySDK(),
		resource.WithDetectors(gcp.NewDetector()),
		resource.WithAttributes(
			// customizable resource attributes
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String(serviceVersion),
		),
	)

	tracerProvider := sdk.NewTracerProvider(
		sdk.WithBatcher(exporter),
		sdk.WithResource(r),
	)
	otel.SetTracerProvider(tracerProvider)

	// setup W3C trace context as global propagator
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tracerProvider
}
