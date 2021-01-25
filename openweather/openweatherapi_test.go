package openweather

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

type params struct {
	city    string
	country string
}

var tests = []struct {
	name     string
	params   params
	expected int
}{
	{"Successful response", params{"Bogota", "co"}, 200},
	{"City not found", params{"", ""}, 404},
	{"Unauthorized", params{"Bogota", "co"}, 401},
	{"Error response", params{"Bogota", "co"}, 503},
}

func newResponder(statusCode int) httpmock.Responder {
	var resp httpmock.Responder

	switch statusCode {
	case http.StatusOK:
		resp, _ = httpmock.NewJsonResponder(200, []byte(`{"coord":{"lon":-0.1257,"lat":51.5085}}`))
	case http.StatusNotFound:
		resp, _ = httpmock.NewJsonResponder(404, []byte(`{"cod":"404","message":"city not found"}`))
	case http.StatusUnauthorized:
		resp, _ = httpmock.NewJsonResponder(401, []byte(`{"cod":401, "message": "Invalid API key. Please see http://openweathermap.org/faq#error401 for more info."}`))
	case http.StatusServiceUnavailable:
		resp = httpmock.NewErrorResponder(fmt.Errorf("Mocked error"))
	}

	return resp
}

func TestGetWeather(t *testing.T) {
	c := NewClient("http://localhost:8081", "1234", "metric").(*clientConfig)
	httpmock.ActivateNonDefault(c.GetClient())
	defer httpmock.DeactivateAndReset()

	for _, test := range tests {
		responder := newResponder(test.expected)
		httpmock.RegisterResponder("GET", "http://localhost:8081/data/2.5/weather", responder)

		t.Run(test.name, func(t *testing.T) {
			statusCode, _ := c.GetWeather(test.params.city, test.params.country)
			if statusCode != test.expected {
				t.Errorf("Error in test: %s\n Got: %v, Expected: %v", test.name, statusCode, test.expected)
			}
		})
	}
}

func TestGetForecast(t *testing.T) {
	c := NewClient("http://localhost:8081", "1234", "metric").(*clientConfig)
	httpmock.ActivateNonDefault(c.GetClient())
	defer httpmock.DeactivateAndReset()

	for _, test := range tests {
		responder := newResponder(test.expected)
		httpmock.RegisterResponder("GET", "http://localhost:8081/data/2.5/forecast", responder)

		t.Run(test.name, func(t *testing.T) {
			statusCode, _ := c.GetForecast(test.params.city, test.params.country)
			if statusCode != test.expected {
				t.Errorf("Error in test: %s\n Got: %v, Expected: %v", test.name, statusCode, test.expected)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	host := "http://localhost:8081"
	c := NewClient(host, "1234", "metric").(*clientConfig)

	if c.HostURL != host {
		t.Errorf("Different hosts. Got: %s, Expected: %s", c.HostURL, host)
	}

	if c.QueryParam.Get("appId") != "1234" {
		t.Errorf("Different apiKeys. Got: %s, Expected: %s", c.QueryParam.Get("appId"), "1234")
	}
}
