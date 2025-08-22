package web

import (
	"net/http"

	"risuhunnik/pkg/database"

	"github.com/labstack/echo/v4"
)

func Index(c echo.Context) error {

	conundrums, err := database.GetVerifiedConundrums()
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "index.html", conundrums)
}
