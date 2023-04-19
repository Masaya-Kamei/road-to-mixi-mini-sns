package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func top(c echo.Context) error {
	return c.String(http.StatusOK, "minimal_sns_app")
}
