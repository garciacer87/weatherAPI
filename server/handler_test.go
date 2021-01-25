package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type params struct {
	city    string
	country string
}

func mockRoute(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

type mockServer struct {
	*gin.Engine
}

func (s *mockServer) makeRequest(city, country string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/test?city=%s&country=%s", city, country), nil)

	s.ServeHTTP(w, req)

	return w
}

type mockService struct{}

func (ms *mockService) GetWeather(city, country string) (int, []byte) {
	if city == "Paris" {
		return 200, nil
	} else if city == "asdfas" {
		return 404, nil
	}

	return 500, nil
}

func TestGetWeather(t *testing.T) {
	mockServer := mockServer{gin.New()}
	mockService := &mockService{}
	mockServer.GET("/test", GetWeather(mockService))

	tests := []struct {
		name     string
		params   params
		expected int
	}{
		{"Successful response", params{"Paris", "fr"}, 200},
		{"Successful response", params{"asdfas", "fr"}, 404},
		{"Successful response", params{"", ""}, 500},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp := mockServer.makeRequest(test.params.city, test.params.country)
			if resp.Code != test.expected {
				t.Errorf("Error in test:  %s. Got: %d, Expected: %d", test.name, resp.Code, test.expected)
			}
		})
	}
}
