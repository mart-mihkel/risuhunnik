package web

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"risuhunnik/pkg/database"

	"github.com/labstack/echo/v4"
)

func CommentForm(c echo.Context) error {

	if !cookiesAgreed(&c) {
		return fmt.Errorf("cookies not agreed!")
	}

	comment := c.FormValue("comment")
	id, err := strconv.Atoi(c.FormValue("conundrum-id"))
	if err != nil {
		return fmt.Errorf("got malfordmed id: %w", err)
	}

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

	err = database.InsertComment(id, comment, author)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "comment-form-result", id)
}
