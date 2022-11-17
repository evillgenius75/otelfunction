package getweather

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type GeoResponse []struct {
	Importance  float64  `json:"importance"`
	Licence     string   `json:"licence"`
	Class       string   `json:"class"`
	OsmID       int      `json:"osm_id"`
	DisplayName string   `json:"display_name"`
	OsmType     string   `json:"osm_type"`
	Lon         string   `json:"lon"`
	PlaceID     int      `json:"place_id"`
	Boundingbox []string `json:"boundingbox"`
	Lat         string   `json:"lat"`
	Type        string   `json:"type"`
}

type WxResponse struct {
	Region            string `json:"region"`
	CurrentConditions struct {
		Dayhour string `json:"dayhour"`
		Temp    struct {
			C int `json:"c"`
			F int `json:"f"`
		} `json:"temp"`
		Precip   string `json:"precip"`
		Humidity string `json:"humidity"`
		Wind     struct {
			Km   int `json:"km"`
			Mile int `json:"mile"`
		} `json:"wind"`
		IconURL string `json:"iconURL"`
		Comment string `json:"comment"`
	} `json:"currentConditions"`
	NextDays []struct {
		Day     string `json:"day"`
		Comment string `json:"comment"`
		MaxTemp struct {
			C int `json:"c"`
			F int `json:"f"`
		} `json:"max_temp"`
		MinTemp struct {
			C int `json:"c"`
			F int `json:"f"`
		} `json:"min_temp"`
		IconURL string `json:"iconURL"`
	} `json:"next_days"`
	ContactAuthor struct {
		Email    string `json:"email"`
		AuthNote string `json:"auth_note"`
	} `json:"contact_author"`
	DataSource string `json:"data_source"`
}

func MakeGeoRequest(ctx context.Context, city, state string) (lat, lon string) {
	// create an instrumented HTTP client
	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	url := fmt.Sprintf("https://forward-reverse-geocoding.p.rapidapi.com/v1/forward?city=%s&state=%s&accept-language=en&polygon_threshold=0.0", city, state)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Add("X-RapidAPI-Key", "8e5d36d1abmsh00ff90acd35e126p108f89jsnd0e08698a22e")
	req.Header.Add("X-RapidAPI-Host", "forward-reverse-geocoding.p.rapidapi.com")
	if err != nil {
		// error handling
		return
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("error connecting to geo API: %v", err)
		return
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	var latLong GeoResponse
	if err := json.Unmarshal(body, &latLong); err != nil {
		fmt.Printf("can not unmarshall JSON error: %v", err)
	}
	// response handling code goes here

	return latLong[0].Lat, latLong[0].Lon
}

func GetWeatherRequest(ctx context.Context, lat, long string) (temp string, err error) {
	//opts := otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request){return operation}("GetWeatherRequest"))
	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	url := fmt.Sprintf("https://weatherdbi.herokuapp.com/data/weather/%s,%s", lat, long)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		fmt.Println("could not access weatherDB api")
		return "", err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("error bad response from weatherDB API: %v", err)
		return "", err
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	var wx WxResponse

	if err := json.Unmarshal(body, &wx); err != nil {
		fmt.Printf("can not unmarshall JSON error: %v", err)
	}

	return strconv.Itoa(wx.CurrentConditions.Temp.F), nil
}
