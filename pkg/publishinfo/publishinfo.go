package publishinfo

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	instrumentationName = "evill-genius.com/my-pubsub-instrumentation-lib"
	instrumentationVer  = "0.1.0"
)

func PublishMessage(ctx context.Context, client *pubsub.Client, msg *pubsub.Message, topicID, projectID string) (string, error) {

	// create span
	ctx, span := beforePublishMessage(ctx, topicID, msg)
	defer span.End()

	// Send Pub/Sub message
	id, err := client.Topic(topicID).Publish(ctx, msg).Get(ctx)
	if err != nil {
		return "", fmt.Errorf("error: %v", err)
	}
	// enrich span with publish result
	afterPublishMessage(span, id, err)
	return id, nil

}

func beforePublishMessage(ctx context.Context, topicID string, msg *pubsub.Message) (context.Context, trace.Span) {
	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindProducer),
		trace.WithAttributes(
			// customizable span attributes
			semconv.MessagingSystemKey.String("pubsub"),
			semconv.MessagingDestinationKey.String(topicID),
			semconv.MessagingDestinationKindTopic,
		),
	}

	tracer := otel.Tracer(
		instrumentationName, trace.WithInstrumentationVersion(instrumentationVer),
	)
	ctx, span := tracer.Start(ctx, fmt.Sprintf("%s send", topicID), opts...)

	if msg.Attributes == nil {
		msg.Attributes = make(map[string]string)
	}

	// propagate Span across process boundaries
	otel.GetTextMapPropagator().Inject(ctx, propagation.MapCarrier(msg.Attributes))

	return ctx, span
}

func afterPublishMessage(span trace.Span, messageID string, err error) {
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetAttributes(semconv.MessagingMessageIDKey.String(messageID))
	}
}
