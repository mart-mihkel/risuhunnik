package web

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"

	"risuhunnik/pkg/database"

	"github.com/labstack/echo/v4"
)

func ToggleStar(c echo.Context) error {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return fmt.Errorf("malfordmed id: %w", err)
	}

	cookie, err := getCookie(&c)
	if err != nil {
		return err
	}

	value, err := deserializeCookie(cookie)
	if err != nil {
		return err
	}

	i := slices.Index(value.Starred, id)
	isStarred := i >= 0
	if isStarred {
		value.Starred = slices.Delete(value.Starred, i, i+1)
	} else {
		value.Starred = append(value.Starred, id)
	}

	escaped, err := serializeCookieValue(value)
	if err != nil {
		return err
	}

	cookie.Value = escaped
	c.SetCookie(cookie)

	var conundrum *database.Conundrum
	if isStarred {
		conundrum, err = database.UnStarConundrum(id)
	} else {
		conundrum, err = database.StarConundrum(id)
	}

	if err != nil {
		return err
	}

	res := &ConundrumResult{
		Conundrum: conundrum,
		Next:      conundrum.Id + 1,
		Prev:      conundrum.Id - 1,
		IsStarred: !isStarred,
	}

	return c.Render(http.StatusOK, "conundrum-stars", res)
}
