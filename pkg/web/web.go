package web

import (
	"fmt"
	"net/http"
	"strconv"

	"risuhunnik/pkg/database"

	"github.com/labstack/echo/v4"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

const searchResultLimit = 5

func Index(c echo.Context) error {
	cs, err := database.GetAllConundrums()
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "index.html", cs)
}

func StarButton(c echo.Context) error {
	sid := c.QueryParam("id")

	id, err := strconv.Atoi(sid)
	if err != nil {
		return fmt.Errorf("got malformed id: %w", err)
	}

	co, err := database.StarConundrum(id)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "star-button", co)
}

func Conundrums(c echo.Context) error {
	cs, err := database.GetAllConundrums()
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "conundrums", cs)
}

func SearchModal(c echo.Context) error {
	cs, err := database.GetAllConundrums()
	if err != nil {
		return err
	}

	s := c.FormValue("search")
	var out []database.Conundrum
	for i, co := range cs {
		if i == searchResultLimit {
			break
		}

		if !fuzzy.MatchFold(s, co.Text) {
			continue
		}

		out = append(out, co)
	}

	return c.Render(http.StatusOK, "search-results", out)
}

func AddModal(c echo.Context) error {
	text := c.FormValue("text")

	co := &database.Conundrum{Text: text}

	id, err := database.InsertConundrum(co)
	if err != nil {
		return c.Render(http.StatusOK, "add-modal", "Upload failed!")
	}

	co.Id = id

	return c.Render(http.StatusOK, "add-modal", "Conundrum uploaded!")
}

func Modal(c echo.Context) error {
	m := c.QueryParam("modal")

	if m == "add" {
		return c.Render(http.StatusOK, "add-modal", nil)
	}

	if m == "search" {
		return c.Render(http.StatusOK, "search-modal", nil)
	}

	return c.Render(http.StatusOK, "hidden-modal", nil)
}
