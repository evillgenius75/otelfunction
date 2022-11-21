// Package helloworld provides a set of Cloud Functions samples.
package helloworld

import (
	//"encoding/json"

	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	getweather "github.com/evillgenius75/otel-function/pkg/get-weather"
	"github.com/evillgenius75/otel-function/pkg/instrumentor"
	"github.com/evillgenius75/otel-function/pkg/otelsetup"
	"github.com/evillgenius75/otel-function/pkg/publishinfo"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

var InstrumentedHandler instrumentor.HttpHandler
var projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
var topicID = os.Getenv("TOPIC_ID")
var psClient *pubsub.Client

func init() {
	tracerProvider := otelsetup.InitTracing("hello-world-svc", "0.1.0")
	InstrumentedHandler = instrumentor.InstrumentedHandler("HelloGet", helloGet, tracerProvider)
	functions.HTTP("HelloGet", InstrumentedHandler)
}

// helloGet is an HTTP Cloud Function.
func helloGet(w http.ResponseWriter, r *http.Request) {
	var err error
	psClient, err = pubsub.NewClient(r.Context(), projectID)
	if err != nil {
		log.Fatalf("pubsub: NewClient: %v", err)
	}
	city := r.URL.Query().Get("city")
	state := r.URL.Query().Get("state")
	fmt.Fprintf(w, "city => %s\tstate => %s\n", city, state)
	lat, long := getweather.MakeGeoRequest(r.Context(), city, state)
	fmt.Fprintf(w, "latitude => %s\tlongitude => %s\n", lat, long)
	temp, err := getweather.GetWeatherRequest(r.Context(), lat, long)
	if err != nil {
		fmt.Fprintf(w, "Error getting weather details: %v\n", err)
	}
	fmt.Fprintf(w, "current temp => %s\n", temp)
	fmt.Fprintf(w, "Sending info to PubSub Topic: %s in Project: %s\n", topicID, projectID)
	fmt.Printf("Sending info to PubSub Topic: %s in Project: %s\n", topicID, projectID)

	//msgJSON, _ := json.Marshal(userInfo)
	//fmt.Fprintf(w, "Marshalled JSON Data: %s", string(msgJSON))
	msg := &pubsub.Message{
		Data: []byte(fmt.Sprintf("city:%s, state:%s, lattitude:%s, longitude:%s, temp:%s, created_at:%s", city, state, lat, long, temp, time.Now().String())),
	}
	fmt.Fprint(w, string(msg.Data))

	id, err := publishinfo.PublishMessage(r.Context(), psClient, msg, topicID, projectID)
	if err != nil {
		fmt.Fprintf(w, "topic(%s).Publish.Get: %v\n", topicID, err)
	} else {
		fmt.Fprintf(w, "Published a message; msg ID: %v\n", id)
	}
}
