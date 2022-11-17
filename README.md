<img src="https://avatars2.githubusercontent.com/u/2810941?v=3&s=96" alt="Google Cloud Platform logo" title="Google Cloud Platform" align="right" height="96" width="96"/>

# Cloud Functions Observability with OpenTelemtry

This prototype is to show the power of the go OpenTelemetry library, Cloud Trace and Cloud Logging to get observability into your event driven cloud functions.

deploy to a cloud function v2 with the preview go1.19 runtime.

You can then call a curl command to 

```
curl -m 70 -X GET "https://function-<instance>-uc.a.run.app?city=Austin&state=TX" -H "Authorization: bearer $(gcloud auth print-identity-token)" -H "Content-Type: application/json" -d ''

```

This will call 2 downstream functions that depend on 2 other APIs 
1. Reverse Geo to get Lat,Long from City,State
2. Get weather data based on Lat, Long and display the Current temp

The Trace will span across the calls and show you which API calls take the longest to complete

Ensure that your Cloud Run Service Account has Cloud Trace IAM permissions to create traces
