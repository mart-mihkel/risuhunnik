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

	err = database.InsertComment(id, comment)
	if err != nil {
		return c.Render(http.StatusOK, "comment-form-result", nil)
	}

	return c.Render(http.StatusOK, "comment-form-result", id)
}
