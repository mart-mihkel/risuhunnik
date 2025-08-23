package web

import (
	"net/http"

	"risuhunnik/pkg/database"

	"github.com/labstack/echo/v4"
)

func UploadForm(c echo.Context) error {

	text := c.FormValue("conundrum")
	cookie, err := maybeInitCookie(&c)
	if err != nil {
		return err
	}

	value, err := deserializeCookie(cookie)
	if err != nil {
		return err
	}

	err = database.InsertConundrum(text, value.Author)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "upload-form-result", true)
}
