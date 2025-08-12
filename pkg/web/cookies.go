package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"risuhunnik/pkg/database"
	"slices"

	"github.com/labstack/echo/v4"
)

type cookieValue struct {
	Starred []int  `json:"starred"`
	Author  string `json:"author"`
}

func initCookie() (*http.Cookie, error) {
	author, err := database.RandomAuthor()
	if err != nil {
		return nil, err
	}

	value := cookieValue{Author: author}
	valuebytes, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize cookie value: %w", err)
	}

	escaped := url.QueryEscape(string(valuebytes))

	return &http.Cookie{
		Name:  "risuhunnik-cookie",
		Value: escaped,
	}, nil
}

func getCookie(c *echo.Context) (*http.Cookie, error) {
	cookie, err := (*c).Cookie("risuhunnik-cookie")
	if err != nil {
		return initCookie()
	}

	return cookie, nil
}

func deserializeCookie(cookie *http.Cookie) (*cookieValue, error) {
	unscaped, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to unscape cookie value: %s", err)
	}

	var value cookieValue
	err = json.Unmarshal([]byte(unscaped), &value)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize cookie value: %w", err)
	}

	return &value, nil
}

func serializeCookieValue(value *cookieValue) (string, error) {
	valuebytes, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("failed to serialize cookie value: %w", err)
	}

	escaped := url.QueryEscape(string(valuebytes))

	return escaped, nil
}

func isStarred(id int, c *echo.Context) (bool, error) {

	cookie, err := (*c).Cookie("risuhunnik-cookie")
	if err != nil {
		return false, nil
	}

	value, err := deserializeCookie(cookie)
	if err != nil {
		return false, err
	}

	return slices.Contains(value.Starred, id), nil
}
