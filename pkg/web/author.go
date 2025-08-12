package web

import (
	"net/http"

	"risuhunnik/pkg/database"

	"github.com/labstack/echo/v4"
)

func Author(c echo.Context) error {

	type Result struct {
		Author     string
		Stars      int
		Conundrums []database.Conundrum
		Comments   []database.Comment
	}

	author := c.QueryParam("author")
	conundrums, err := database.GetAuthorConundrums(author)
	if err != nil {
		return err
	}

	stars, err := database.GetAuthorStars(author)
	if err != nil {
		return err
	}

	comments, err := database.GetAuthorComments(author)
	if err != nil {
		return err
	}

	res := &Result{
		Author:     author,
		Stars:      stars,
		Conundrums: conundrums,
		Comments:   comments,
	}

	return c.Render(http.StatusOK, "author", res)
}
