package web

import (
	"fmt"
	"net/http"
	"strconv"

	"risuhunnik/pkg/database"

	"github.com/labstack/echo/v4"
)

type ConundrumsResult struct {
	Conundrums    []database.Conundrum
	CookiesAgreed bool
}

type ConundrumResult struct {
	Conundrum     *database.Conundrum
	Comments      []database.Comment
	Next          int
	Prev          int
	IsStarred     bool
	CookiesAgreed bool
}

func Index(c echo.Context) error {
	conundrums, err := database.GetAllConundrums()
	if err != nil {
		return err
	}

	res := &ConundrumsResult{
		Conundrums:    conundrums,
		CookiesAgreed: cookiesAgreed(&c),
	}

	return c.Render(http.StatusOK, "index.html", res)
}

func Cookies(c echo.Context) error {
	agreed, err := strconv.ParseBool(c.QueryParam("agreed"))
	if err != nil {
		return fmt.Errorf("got malfordmed agreed: %w", err)
	}

	if !agreed {
		return c.NoContent(http.StatusOK)
	}

	author, err := database.RandomAuthor()
	if err != nil {
		return err
	}

	c.SetCookie(&http.Cookie{Name: "agreed"})
	c.SetCookie(&http.Cookie{Name: "author", Value: author})

	return c.NoContent(http.StatusOK)
}

func Conundrums(c echo.Context) error {
	conundrums, err := database.GetAllConundrums()
	if err != nil {
		return err
	}

	res := &ConundrumsResult{
		Conundrums:    conundrums,
		CookiesAgreed: cookiesAgreed(&c),
	}

	return c.Render(http.StatusOK, "conundrums", res)
}

func Conundrum(c echo.Context) error {

	id, err := strconv.Atoi(c.QueryParam("id"))
	if err != nil {
		return fmt.Errorf("got malfordmed id: %w", err)
	}

	conundrum, err := database.GetConundrum(id)
	if err != nil {
		return err
	}

	comments, err := database.GetConundrumComments(id)
	if err != nil {
		return err
	}

	starred, err := isStarred(id, &c)
	if err != nil {
		return err
	}

	res := &ConundrumResult{
		Conundrum:     conundrum,
		Comments:      comments,
		Next:          conundrum.Id + 1,
		Prev:          conundrum.Id - 1,
		IsStarred:     starred,
		CookiesAgreed: cookiesAgreed(&c),
	}

	return c.Render(http.StatusOK, "conundrum", res)
}

func cookiesAgreed(c *echo.Context) bool {
	_, err := (*c).Cookie("agreed")
	return err == nil
}
