package web

import (
	"net/http"

	"risuhunnik/pkg/database"

	"github.com/labstack/echo/v4"
)

func UploadForm(c echo.Context) error {

	if !cookiesAgreed(&c) {
		// TODO: render a response
		return c.Render(http.StatusOK, "comment-form-result", nil)
	}

	text := c.FormValue("conundrum")

	err := database.InsertConundrum(text)
	if err != nil {
		return c.Render(http.StatusOK, "upload-form-result", nil)
	}

	return c.Render(http.StatusOK, "upload-form-result", true)
}
