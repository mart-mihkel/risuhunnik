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

func Star(c echo.Context) error {

	type Result struct {
		Conundrum     *database.Conundrum
		IsStarred     bool
		CookiesAgreed bool
	}

	if !cookiesAgreed(&c) {
		// TODO: render a response
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

	if slices.Contains(value, id) {
		// TODO: render a response
		return fmt.Errorf("conundrum already starred!")
	}

	value = append(value, id)
	jsonbytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to serialize cookie: %w", err)
	}

	c.SetCookie(&http.Cookie{Name: "starred", Value: string(jsonbytes)})

	conundrum, err := database.StarConundrum(id)
	if err != nil {
		return err
	}

	res := &Result{
		Conundrum:     conundrum,
		IsStarred:     true,
		CookiesAgreed: cookiesAgreed(&c),
	}

	return c.Render(http.StatusOK, "conundrum-stars", res)
}

func IsStarred(id int, c *echo.Context) (bool, error) {

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
