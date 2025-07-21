package pages

import (
	"fmt"
	"risuhunnik/pkg/database"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func Index(c echo.Context) error {
	cs, err := database.GetAllConundrums()
	if err != nil {
		return err
	}

	return c.Render(200, "index.html", cs)
}

func Star(c echo.Context) error {
	sid := c.QueryParam("id")

	id, err := strconv.Atoi(sid)
	if err != nil {
		return fmt.Errorf("got malformed id: %w ", err)
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

	return c.Render(200, "button-star", co)
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

	return c.Render(200, "tags", counts)
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

	return c.Render(200, "conundrums", cs)
}

func SearchConundrums(c echo.Context) error {
	s := c.FormValue("search")

	cs, err := database.GetConundrumsBySubstring(s)
	if err != nil {
		return err
	}

	return c.Render(200, "conundrums", cs)
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

	return c.Render(200, "modal-add", co)
}

func Modal(c echo.Context) error {
	m := c.QueryParam("modal")

	if m == "add" {
		return c.Render(200, "modal-add", nil)
	}

	if m == "search" {
		return c.Render(200, "modal-search", nil)
	}

	return c.Render(200, "modal-hidden", nil)
}
