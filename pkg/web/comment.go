package web

import (
	"fmt"
	"net/http"
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

		cookie = &http.Cookie{Name: "author", Value: author}
		c.SetCookie(cookie)
	}

	err = database.InsertComment(id, comment, cookie.Value)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "comment-form-result", id)
}
