package openweather

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

//Client used to make requests to openweathermap.org API
type Client interface {
	GetWeather(city, country string) (int, []byte)
	GetForecast(city, country string) (int, []byte)
}

//clientConfig struct used to store config attributes necessary to connect to openweathermap.org API
type clientConfig struct {
	*resty.Client
}

//NewClient retrieves a new OpenWheater client
func NewClient(host, apiKey, unit string) Client {
	c := &clientConfig{resty.New()}

	c.SetHostURL(host).
		SetRetryCount(3).
		SetQueryParam("appid", apiKey).
		SetQueryParam("units", unit)

	return c
}

//GetWeather makes a GET request to openweather client to get weather info for a specific city
func (c *clientConfig) GetWeather(city, country string) (int, []byte) {
	resp, err := c.R().
		SetQueryParams(map[string]string{
			"q": fmt.Sprintf("%s,%s", city, country),
		}).Get("/data/2.5/weather")

	if err != nil {
		return http.StatusServiceUnavailable, []byte(`{"code":503, "message":"Error making request to OpenWeather API"`)
	}

	return resp.StatusCode(), resp.Body()
}

//GetForecast makes a GET request to openweather client to get forecast info for a specific city
func (c *clientConfig) GetForecast(city, country string) (int, []byte) {
	resp, err := c.R().
		SetQueryParams(map[string]string{
			"q":   fmt.Sprintf("%s,%s", city, country),
			"cnt": "3",
		}).Get("/data/2.5/forecast")

	if err != nil {
		return http.StatusServiceUnavailable, []byte(`{"code":503, "message":"Error making request to OpenWeather API"`)
	}

	return resp.StatusCode(), resp.Body()
}
