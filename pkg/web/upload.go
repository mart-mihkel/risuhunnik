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
	cookie, err := c.Cookie("author")
	if err != nil {
		author, err := database.RandomAuthor()
		if err != nil {
			return err
		}

		cookie = &http.Cookie{Name: "author", Value: author}
		c.SetCookie(cookie)
	}

	err = database.InsertConundrum(text, cookie.Value)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "upload-form-result", true)
}
