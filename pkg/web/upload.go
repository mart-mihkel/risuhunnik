package web

import (
	"fmt"
	"net/http"
	"net/url"

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

		author = url.QueryEscape(author)

		cookie = &http.Cookie{Name: "author", Value: author}
		c.SetCookie(cookie)
	}

	author, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return fmt.Errorf("failed to query unescape author: %w", err)
	}

	err = database.InsertConundrum(text, author)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "upload-form-result", true)
}
