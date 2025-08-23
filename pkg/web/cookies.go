package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"slices"

	"risuhunnik/pkg/database"

	"github.com/labstack/echo/v4"
)

type cookieValue struct {
	Starred []int  `json:"starred"`
	Author  string `json:"author"`
	Token   string `json:"token"`
}

func initCookie(c *echo.Context) (*http.Cookie, error) {
	author, err := database.RandomAuthor()
	if err != nil {
		return nil, err
	}

	value := cookieValue{Author: author}
	valuebytes, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("serializing cookie value: %w", err)
	}

	escaped := url.QueryEscape(string(valuebytes))
	cookie := &http.Cookie{
		Name:  "risuhunnik-cookie",
		Value: escaped,
		Path:  "/",
	}

	(*c).SetCookie(cookie)

	return cookie, nil
}

func maybeInitCookie(c *echo.Context) (*http.Cookie, error) {

	cookie, err := (*c).Cookie("risuhunnik-cookie")
	if err != nil {
		return initCookie(c)
	}

	return cookie, nil
}

func deserializeCookie(cookie *http.Cookie) (*cookieValue, error) {
	unscaped, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return nil, fmt.Errorf("unscaping cookie value: %s", err)
	}

	var value cookieValue
	err = json.Unmarshal([]byte(unscaped), &value)
	if err != nil {
		return nil, fmt.Errorf("deserializing cookie value: %w", err)
	}

	return &value, nil
}

func serializeCookieValue(value *cookieValue) (string, error) {
	valuebytes, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("serializing cookie value: %w", err)
	}

	escaped := url.QueryEscape(string(valuebytes))

	return escaped, nil
}

func hasValidToken(c *echo.Context) (bool, error) {

	cookie, err := maybeInitCookie(c)
	if err != nil {
		return false, err
	}

	value, err := deserializeCookie(cookie)
	if err != nil {
		return false, err
	}

	return database.CheckToken(value.Token)
}

func isStarred(id int, c *echo.Context) (bool, error) {

	cookie, err := maybeInitCookie(c)
	if err != nil {
		return false, err
	}

	value, err := deserializeCookie(cookie)
	if err != nil {
		return false, err
	}

	return slices.Contains(value.Starred, id), nil
}
