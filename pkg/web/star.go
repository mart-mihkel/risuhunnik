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

	cookie, err := maybeInitCookie(&c)
	if err != nil {
		return err
	}

	value, err := deserializeCookie(cookie)
	if err != nil {
		return err
	}

	i := slices.Index(value.Starred, id)
	starred := i >= 0
	if starred {
		value.Starred = slices.Delete(value.Starred, i, i+1)
	} else {
		value.Starred = append(value.Starred, id)
	}

	escaped, err := serializeCookieValue(value)
	if err != nil {
		return err
	}

	cookie.Value = escaped
	cookie.Path = "/"

	c.SetCookie(cookie)

	var conundrum *database.Conundrum
	if starred {
		conundrum, err = database.UnStarConundrum(id)
	} else {
		conundrum, err = database.StarConundrum(id)
	}

	if err != nil {
		return err
	}

	res := &ConundrumResult{
		Conundrum: conundrum,
		Starred:   !starred,
	}

	return c.Render(http.StatusOK, "conundrum-stars", res)
}
