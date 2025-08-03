package web

import (
	"encoding/json"
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

	cookie, err := c.Cookie("starred")
	if err != nil {
		cookie = &http.Cookie{Name: "starred", Value: "[]"}
	}

	var value []int
	err = json.Unmarshal([]byte(cookie.Value), &value)
	if err != nil {
		return fmt.Errorf("failed to deserialize cookie: %w", err)
	}

	i := slices.Index(value, id)
	starred := i >= 0
	if starred {
		value = slices.Delete(value, i, i+1)
	} else {
		value = append(value, id)
	}

	jsonbytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to serialize cookie: %w", err)
	}

	c.SetCookie(&http.Cookie{Name: "starred", Value: string(jsonbytes)})

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
		Conundrum:     conundrum,
		IsStarred:     !starred,
		CookiesAgreed: cookiesAgreed(&c),
	}

	return c.Render(http.StatusOK, "conundrum-stars", res)
}

func isStarred(id int, c *echo.Context) (bool, error) {

	cookie, err := (*c).Cookie("starred")
	if err != nil {
		cookie = &http.Cookie{Name: "starred", Value: "[]"}
	}

	var value []int
	err = json.Unmarshal([]byte(cookie.Value), &value)
	if err != nil {
		return false, fmt.Errorf("failed to deserialize cookie: %w", err)
	}

	return slices.Contains(value, id), nil
}
