package web

import (
	"fmt"
	"net/http"

	"risuhunnik/pkg/database"

	"github.com/labstack/echo/v4"
)

func UploadForm(c echo.Context) error {

	if !cookiesAgreed(&c) {
		return fmt.Errorf("cookies not agreed!")
	}

	text := c.FormValue("conundrum")
	cookie, err := getOrMakeCookie(&c)
	if err != nil {
		return err
	}

	c.SetCookie(cookie)

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
