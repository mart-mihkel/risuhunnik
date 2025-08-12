package web

import (
	"fmt"
	"net/http"
	"strconv"

	"risuhunnik/pkg/database"

	"github.com/labstack/echo/v4"
)

func CommentForm(c echo.Context) error {

	comment := c.FormValue("comment")
	id, err := strconv.Atoi(c.FormValue("conundrum-id"))
	if err != nil {
		return fmt.Errorf("got malfordmed id: %w", err)
	}

	cookie, err := getCookie(&c)
	if err != nil {
		return err
	}

	c.SetCookie(cookie)

	value, err := deserializeCookie(cookie)
	if err != nil {
		return err
	}

	err = database.InsertComment(id, comment, value.Author)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "comment-form-result", id)
}
