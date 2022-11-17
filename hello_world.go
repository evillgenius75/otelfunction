// Package helloworld provides a set of Cloud Functions samples.
package helloworld

import (
	"fmt"
	"net/http"

	getweather "github.com/evillgenius75/otel-function/pkg/get-weather"
	"github.com/evillgenius75/otel-function/pkg/instrumentor"
	"github.com/evillgenius75/otel-function/pkg/otelsetup"
	"go.opentelemetry.io/otel/baggage"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

var InstrumentedHandler instrumentor.HttpHandler

func init() {
	tracerProvider := otelsetup.InitTracing("hello-world-svc", "0.1.0")
	InstrumentedHandler = instrumentor.InstrumentedHandler("HelloGet", helloGet, tracerProvider)
	functions.HTTP("HelloGet", InstrumentedHandler)
}

// helloGet is an HTTP Cloud Function.
func helloGet(w http.ResponseWriter, r *http.Request) {
	bag, _ := baggage.New()
	newField, _ := baggage.NewMember("functionName", "getWeatherRequest")
	bag.SetMember(newField)
	ctx := baggage.ContextWithBaggage(r.Context(), bag)
	city := r.URL.Query().Get("city")
	state := r.URL.Query().Get("state")
	fmt.Fprintf(w, "city => %s\tstate => %s\n", city, state)
	lat, long := getweather.MakeGeoRequest(r.Context(), city, state)
	fmt.Fprintf(w, "latitude => %s\tlongitude => %s\n", lat, long)
	temp, err := getweather.GetWeatherRequest(ctx, lat, long)
	if err != nil {
		fmt.Fprintf(w, "Error getting weather details: %v", err)
	}
	fmt.Fprintf(w, "current temp => %s\n", temp)
}
