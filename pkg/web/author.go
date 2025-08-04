package web

import (
	"net/http"

	"risuhunnik/pkg/database"

	"github.com/labstack/echo/v4"
)

func Author(c echo.Context) error {

	type Result struct {
		Author     string
		Conundrums []database.Conundrum
	}

	author := c.QueryParam("author")
	conundrums, err := database.GetAuthorConundrums(author)
	if err != nil {
		return err
	}

	res := &Result{
		Author: author,
		Conundrums: conundrums,
	}

	return c.Render(http.StatusOK, "author", res)
}
