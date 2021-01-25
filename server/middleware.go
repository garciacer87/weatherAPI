package server

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

//ValidateRequest returns a handler used as middleware to validate query params from incoming requests
func ValidateRequest() gin.HandlerFunc {
	cityRexp, _ := regexp.Compile(`^[a-zA-Z\s]+$`)
	countryRexp, _ := regexp.Compile(`^[a-z]{2}$`)

	return func(c *gin.Context) {
		errors := make([]string, 0)

		for _, param := range []string{"city", "country"} {
			value, ok := c.GetQuery(param)
			if !ok {
				errors = append(errors, fmt.Sprintf("missing query param: '%s'", param))
				continue
			}
			if value == "" {
				errors = append(errors, fmt.Sprintf("%s cannot be empty", param))
			}
		}

		if !cityRexp.MatchString(c.Query("city")) {
			errors = append(errors, "city must be a string")
		}

		if !countryRexp.MatchString(c.Query("country")) {
			errors = append(errors, "country must be a two characters string in lowercase")
		}

		if len(errors) > 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": errors})
		}
	}
}
