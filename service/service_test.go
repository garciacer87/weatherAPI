package service

import (
	"encoding/json"
	"testing"
)

var (
	weatherResp  = []byte(`{"coord":{"lon":2.3488,"lat":48.8534},"weather":[{"id":803,"main":"Clouds","description":"broken clouds","icon":"04n"}],"base":"stations","main":{"temp":1.85,"feels_like":-5.05,"temp_min":1,"temp_max":2.22,"pressure":1002,"humidity":93},"visibility":10000,"wind":{"speed":7.2,"deg":290},"clouds":{"all":75},"dt":1611558107,"sys":{"type":1,"id":6550,"country":"FR","sunrise":1611559763,"sunset":1611592574},"timezone":3600,"id":2988507,"name":"Paris","cod":200}`)
	forecastResp = []byte(`{"cod":"200","message":0,"cnt":2,"list":[{"dt":1611565200,"main":{"temp":2.27,"feels_like":-3.25,"temp_min":2.27,"temp_max":2.71,"pressure":1004,"sea_level":1004,"grnd_level":1001,"humidity":87,"temp_kf":-0.44},"weather":[{"id":803,"main":"Clouds","description":"broken clouds","icon":"04d"}],"clouds":{"all":71},"wind":{"speed":5.12,"deg":336},"visibility":10000,"pop":0,"sys":{"pod":"d"},"dt_txt":"2021-01-25 09:00:00"},{"dt":1611576000,"main":{"temp":4.1,"feels_like":-1.63,"temp_min":4.1,"temp_max":4.72,"pressure":1007,"sea_level":1007,"grnd_level":1004,"humidity":74,"temp_kf":-0.62},"weather":[{"id":803,"main":"Clouds","description":"broken clouds","icon":"04d"}],"clouds":{"all":70},"wind":{"speed":5.33,"deg":343},"visibility":10000,"pop":0,"sys":{"pod":"d"},"dt_txt":"2021-01-25 12:00:00"}],"city":{"id":2988507,"name":"Paris","coord":{"lat":48.8534,"lon":2.3488},"country":"FR","population":2138551,"timezone":3600,"sunrise":1611559763,"sunset":1611592574}}`)
)

type mockService struct{}

func (ms *mockService) GetWeather(city, country string) (int, []byte) {
	if city == "Paris" {
		return 200, weatherResp
	} else if city == "asdf" {
		return 404, nil
	}

	return 200, nil
}

func (ms *mockService) GetForecast(city, country string) (int, []byte) {
	if city == "Paris" {
		return 200, forecastResp
	} else if city == "qwer" {
		return 404, nil
	}

	return 200, nil
}

type mockCache struct {
	v map[string][]byte
}

func (mc *mockCache) SetValue(id string, v []byte) {
	mc.v[id] = v
}

func (mc *mockCache) GetValue(id string) []byte {
	return mc.v[id]
}

type params struct {
	city    string
	country string
}

func TestGetWeather(t *testing.T) {
	s := New("host", "apikey", "metric", 2)

	ms := s.(*service)
	ms.apiClient = &mockService{}
	ms.cache = &mockCache{make(map[string][]byte)}

	tests := []struct {
		name     string
		params   params
		expected int
	}{
		{"Succesful response", params{"Paris", "FR"}, 200},
		{"Failed weather response", params{"asdf", "zz"}, 404},
		{"Failed forecast response", params{"qwer", "zz"}, 404},
		{"Successful response from cache", params{"Paris", "FR"}, 200},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			statusCode, _ := s.GetWeather(test.params.city, test.params.country)
			if statusCode != test.expected {
				t.Errorf("Error in test:  %s. Got: %d, Expected: %d", test.name, statusCode, test.expected)
			}
		})
	}
}

func TestRespBuilder(t *testing.T) {
	var resp Response
	unit := "metric"

	data, _ := buildResponse(weatherResp, forecastResp, unit)

	json.Unmarshal(data, &resp)

	if resp.Cloudiness != "broken clouds" {
		t.Errorf("Error in cloudiness: Got: %s, Expected: %s", resp.Cloudiness, "broken clouds")
	}

	if resp.Temp != "2ºC" {
		t.Errorf("Error in temperature: Got: %s, Expected: %s", resp.Temp, "2ºC")
	}

	if resp.Location != "Paris, FR" {
		t.Errorf("Error in temperature: Got: %s, Expected: %s", resp.Location, "Paris, FR")
	}

	if len(resp.Forecast) != 2 {
		t.Errorf("Error in forecast list size: Got: %d, Expected: %d", len(resp.Forecast), 2)
	}

	_, err := buildResponse(weatherResp, []byte(""), unit)
	if err == nil {
		t.Errorf("Expected error ")
	}

	_, err = buildResponse([]byte(""), forecastResp, unit)
	if err == nil {
		t.Errorf("Expected error ")
	}
}

func TestGetRequestId(t *testing.T) {
	id := getRequestID("Paris", "FR")
	if id != "paris_fr" {
		t.Errorf("Id is different than expected. Got: %s, Expected: %s", id, "paris_fr")
	}
}

func TestFmtTemperature(t *testing.T) {
	temp := fmtTemperature(110.25, "metric")
	if temp != "110ºC" {
		t.Errorf("Temperature is different than expected. Got: %s, Expected: %s", temp, "110ºC")
	}
}

func TestFmtTime(t *testing.T) {
	time := fmtTime(1611558107)
	if time != "04:01" {
		t.Errorf("Time is different than expected. Got: %s, Expected: %s", time, "04:01")
	}
}

func TestFmtDateTime(t *testing.T) {
	time := fmtDateTime(1611558107)
	if time != "25/01/2021 04:01" {
		t.Errorf("Time is different than expected. Got: %s, Expected: %s", time, "25/01/2021 04:01")
	}
}

func TestGetWindDirection(t *testing.T) {
	tests := []struct {
		name     string
		deg      int
		expected string
	}{
		{"North direction", 0, "North"},
		{"North-NorthEast direction", 35, "North-NorthEast"},
		{"East-NorthEast direction", 80, "East-NorthEast"},
		{"East-SouthEast direction", 100, "East-SouthEast"},
		{"South-SouthEast direction", 140, "South-SouthEast"},
		{"South-SouthWest direction", 200, "South-SouthWest"},
		{"West-SouthWest direction", 233, "West-SouthWest"},
		{"West-NorthWest direction", 300, "West-NorthWest"},
		{"North-NorthWest direction", 345, "North-NorthWest"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			direction := getWindDirection(test.deg)
			if direction != test.expected {
				t.Errorf("Error in test:  %s. Got: %s, Expected: %s", test.name, direction, test.expected)
			}
		})
	}
}
