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

	if !cookiesAgreed(&c) {
		return fmt.Errorf("cookies not agreed!")
	}

	id, err := strconv.Atoi(c.FormValue("id"))
	if err != nil {
		return fmt.Errorf("got malfordmed id: %w", err)
	}

	cookie, err := getOrMakeCookie(&c)
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
		Conundrum:     conundrum,
		IsStarred:     !isStarred,
		CookiesAgreed: cookiesAgreed(&c),
	}

	return c.Render(http.StatusOK, "conundrum-stars", res)
}
