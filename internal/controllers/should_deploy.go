package controllers

import (
	"math/rand"

	"github.com/labstack/echo/v4"
)

func ShouldDeploy(c echo.Context) error {
	n := rand.Intn(10) + 1

	if n%2 == 0 {
		return c.JSON(200, map[string]interface{}{
			"shouldDeploy": true,
			"message":      "You shall pass!",
		})
	}

	return c.JSON(403, map[string]interface{}{
		"shouldDeploy": false,
		"message":      "You shall not pass!",
	})
}
