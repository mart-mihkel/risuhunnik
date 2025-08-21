package web

import (
	"fmt"
	"net/http"
	"risuhunnik/pkg/database"
	"strconv"

	"github.com/labstack/echo/v4"
)

func Verify(c echo.Context) error {

	valid, err := hasValidToken(&c)
	if err != nil {
		return err
	}

	if !valid {
		return c.NoContent(http.StatusUnauthorized)
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return fmt.Errorf("malfordmed id: %w", err)
	}

	conundrum, err := database.ToggleVerifyConundrum(id)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "conundrum-verified", conundrum)
}
