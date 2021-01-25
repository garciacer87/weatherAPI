package server

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestValidateHeader(t *testing.T) {
	s := mockServer{gin.New()}
	s.Use(ValidateRequest()).GET("/test")

	tests := []struct {
		name     string
		params   params
		expected int
	}{
		{"Successful response", params{"Paris", "fr"}, 200},
		{"Missing city", params{"", "fr"}, 400},
		{"Missing country", params{"Paris", ""}, 400},
		{"Wrong chars on city", params{"P@r1s", "fr"}, 400},
		{"Wrong chars on country", params{"Paris", "f1"}, 400},
		{"Uppercase country", params{"Paris", "FR"}, 400},
		{"Country longer than 2 characters", params{"Paris", "fra"}, 400},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp := s.makeRequest(test.params.city, test.params.country)
			if resp.Code != test.expected {
				t.Errorf("Error in test:  %s. Got: %d, Expected: %d", test.name, resp.Code, test.expected)
			}
		})
	}

}
