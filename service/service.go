package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/garciacer87/weatherAPI/apicache"
	"github.com/garciacer87/weatherAPI/openweather"
)

var (
	units = map[string]struct {
		temp  string
		speed string
	}{
		"metric":   {"ºC", "m/s"},
		"imperial": {"ºF", "miles/hr"},
	}

	windDirections = map[int]string{
		0:   "North",
		45:  "NorthEast",
		90:  "East",
		135: "SouthEast",
		180: "South",
		225: "SouthWest",
		270: "West",
		315: "NorthWest",
	}
)

//Service interface used to implement "get weather" logic
type Service interface {
	GetWeather(city, country string) (int, []byte)
}

type service struct {
	apiClient openweather.Client
	unit      string
	cache     apicache.Cache
}

//New returns a new Service
func New(host, apiKey, unit string, cacheDuration int) Service {
	apiClient := openweather.NewClient(host, apiKey, unit)
	cache := apicache.New(cacheDuration)

	return &service{apiClient, unit, cache}
}

//GetWeather gets weather information from a city. Uses a cache for retrieving response
func (s *service) GetWeather(city, country string) (int, []byte) {
	reqID := getRequestID(city, country)

	finalResp := s.cache.GetValue(reqID)
	if finalResp != nil {
		return http.StatusOK, finalResp
	}

	respCode, weatherBody := s.apiClient.GetWeather(city, country)
	if respCode != http.StatusOK {
		return respCode, weatherBody
	}

	respCode, forecastBody := s.apiClient.GetForecast(city, country)
	if respCode != http.StatusOK {
		return respCode, forecastBody
	}

	finalResp, err := buildResponse(weatherBody, forecastBody, s.unit)
	if err != nil {
		return http.StatusInternalServerError, []byte(`{"code":500, "message":"Error processing response"`)
	}

	s.cache.SetValue(reqID, finalResp)

	return respCode, finalResp
}

func buildResponse(weatherBody, forecastBody []byte, unit string) ([]byte, error) {
	var wResp weatherResponse
	var fcResp forecastResponse

	err := json.Unmarshal(weatherBody, &wResp)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(forecastBody, &fcResp)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	forecastList := make([]forecast, 0)
	for _, fcInfo := range fcResp.Forecast {
		forecastList = append(forecastList, forecast{
			ForecastedDate: fmtDateTime(fcInfo.Dt),
			Temp:           fmtTemperature(fcInfo.Main.Temp, unit),
			Feel:           fmtTemperature(fcInfo.Main.FeelsLike, unit),
			Min:            fmtTemperature(fcInfo.Main.TempMin, unit),
			Max:            fmtTemperature(fcInfo.Main.TempMax, unit),
			Cloudiness:     fcInfo.Weather[0].Description,
			Humidity:       fmt.Sprintf("%v%%", fcInfo.Main.Humidity),
		})
	}

	r := Response{
		Location:   fmt.Sprintf("%s, %s", wResp.Name, wResp.Sys.Country),
		Temp:       fmtTemperature(wResp.Main.Temp, unit),
		Feel:       fmtTemperature(wResp.Main.FeelsLike, unit),
		Min:        fmtTemperature(wResp.Main.TempMin, unit),
		Max:        fmtTemperature(wResp.Main.TempMax, unit),
		Wind:       fmt.Sprintf("%.2f %s %s", wResp.Wind.Speed, units[unit].speed, getWindDirection(wResp.Wind.Deg)),
		Cloudiness: wResp.Weather[0].Description,
		Pressure:   fmt.Sprintf("%v hpa", wResp.Main.Pressure),
		Humidity:   fmt.Sprintf("%v%%", wResp.Main.Humidity),
		Sunrise:    fmtTime(wResp.Sys.Sunrise),
		Sunset:     fmtTime(wResp.Sys.Sunset),
		Coord:      fmt.Sprintf("[%f, %f]", wResp.Coord.Lat, wResp.Coord.Lon),
		ReqTime:    fmt.Sprintf("%02d:%02d", now.Hour(), now.Minute()),
		Forecast:   forecastList,
	}

	finalResp, _ := json.Marshal(&r)

	return finalResp, nil
}

func getRequestID(city, country string) string {
	return strings.ToLower(fmt.Sprintf("%s_%s", city, country))
}

func fmtTemperature(temp float64, unit string) string {
	return fmt.Sprintf("%.0f%s", temp, units[unit].temp)
}

func fmtTime(timestamp int) string {
	datetime := time.Unix(int64(timestamp), 0)
	return fmt.Sprintf("%02d:%02d", datetime.Hour(), datetime.Minute())
}

func fmtDateTime(timestamp int) string {
	datetime := time.Unix(int64(timestamp), 0)
	return fmt.Sprintf("%02d/%02d/%02d %02d:%02d", datetime.Day(), datetime.Month(), datetime.Year(), datetime.Hour(), datetime.Minute())
}

func getWindDirection(deg int) string {
	dir, ok := windDirections[deg]
	if ok {
		return dir
	}

	if 1 < deg && deg < 44 {
		return "North-NorthEast"
	} else if 46 < deg && deg < 89 {
		return "East-NorthEast"
	} else if 91 < deg && deg < 134 {
		return "East-SouthEast"
	} else if 136 < deg && deg < 179 {
		return "South-SouthEast"
	} else if 181 < deg && deg < 224 {
		return "South-SouthWest"
	} else if 226 < deg && deg < 269 {
		return "West-SouthWest"
	} else if 271 < deg && deg < 314 {
		return "West-NorthWest"
	} else if 316 < deg && deg < 359 {
		return "North-NorthWest"
	}

	return ""
}
