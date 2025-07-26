package web

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

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

func Tags(c echo.Context) error {
	cs, err := database.GetAllConundrums()
	if err != nil {
		return err
	}

	counts := make(map[string]int)
	for _, co := range cs {
		for _, t := range co.Tags {
			counts[t]++
		}
	}

	return c.Render(http.StatusOK, "tags", counts)
}

func Conundrums(c echo.Context) error {
	t := c.QueryParam("tag")

	var cs []database.Conundrum
	var err error

	if t != "" {
		cs, err = database.GetConundrumsByTag(t)
	} else {
		cs, err = database.GetAllConundrums()
	}

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

		m := fuzzy.MatchFold(s, co.Text)
		ms := fuzzy.FindFold(s, co.Tags)
		if !m && len(ms) == 0 {
			continue
		}

		out = append(out, co)
	}

	return c.Render(http.StatusOK, "search-results", out)
}

func AddModal(c echo.Context) error {
	text := c.FormValue("text")
	tags := c.FormValue("tags")

	co := &database.Conundrum{
		Text: text,
		Tags: strings.Split(tags, " "),
	}

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
