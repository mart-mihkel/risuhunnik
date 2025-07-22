package pages

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"risuhunnik/pkg/database"

	"github.com/labstack/echo/v4"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

func Index(c echo.Context) error {
	cs, err := database.GetAllConundrums()
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "index.html", cs)
}

func Star(c echo.Context) error {
	sid := c.QueryParam("id")

	id, err := strconv.Atoi(sid)
	if err != nil {
		return fmt.Errorf("got malformed id: %w", err)
	}

	co, err := database.GetConundrum(id)
	if err != nil {
		return err
	}

	co.Stars++

	err = database.UpdateConundrum(co)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "button-star", co)
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

func SearchConundrums(c echo.Context) error {
	cs, err := database.GetAllConundrums()
	if err != nil {
		return err
	}

	lim := 5
	s := c.FormValue("search")
	var out []database.Conundrum
	for _, co := range cs {
		m := fuzzy.MatchFold(s, co.Text)
		ms := fuzzy.FindFold(s, co.Tags)
		if !m && len(ms) == 0 || len(out) == lim {
			continue
		}

		out = append(out, co)
	}

	return c.Render(http.StatusOK, "search-results", out)
}

func PostConundrum(c echo.Context) error {
	text := c.FormValue("text")
	tags := c.FormValue("tags")

	co := &database.Conundrum{
		Text:     text,
		Tags:     strings.Split(tags, " "),
		Verified: false,
		Stars:    0,
	}

	id, err := database.InsertConundrum(co)
	if err != nil {
		return err
	}

	co.Id = id

	return c.Render(http.StatusOK, "modal-add", co)
}

func Modal(c echo.Context) error {
	m := c.QueryParam("modal")

	if m == "add" {
		return c.Render(http.StatusOK, "modal-add", nil)
	}

	if m == "search" {
		return c.Render(http.StatusOK, "modal-search", nil)
	}

	return c.Render(http.StatusOK, "modal-hidden", nil)
}
