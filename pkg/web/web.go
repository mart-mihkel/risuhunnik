package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"risuhunnik/pkg/database"

	"github.com/labstack/echo/v4"
)

type RisuhunnikCookie struct {
	Starred  []int `json:"starred"`
	Posts    []int `json:"posts"`
	Comments []int `json:"comments"`
}

type ConundrumsResult struct {
	Conundrums    []database.Conundrum
	CookiesAgreed bool
}

type ConundrumResult struct {
	Conundrum     *database.Conundrum
	Comments      []database.Comment
	Next          int
	Prev          int
	CookiesAgreed bool
}

func Index(c echo.Context) error {
	cs, err := database.GetAllConundrums()
	if err != nil {
		return err
	}

	res := &ConundrumsResult{
		Conundrums:    cs,
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

	bytes, err := json.Marshal(&RisuhunnikCookie{
		Starred:  []int{},
		Posts:    []int{},
		Comments: []int{},
	})

	if err != nil {
		return fmt.Errorf("failed to serialize default cookie: %w", err)
	}

	c.SetCookie(&http.Cookie{
		Name:  "risuhunnik-cookie",
		Value: string(bytes),
	})

	return c.NoContent(http.StatusOK)
}

func Conundrums(c echo.Context) error {
	cs, err := database.GetAllConundrums()
	if err != nil {
		return err
	}

	res := &ConundrumsResult{
		Conundrums:    cs,
		CookiesAgreed: cookiesAgreed(&c),
	}

	return c.Render(http.StatusOK, "conundrums", res)
}

func Conundrum(c echo.Context) error {
	id, err := strconv.Atoi(c.QueryParam("id"))
	if err != nil {
		return fmt.Errorf("got malfordmed id: %w", err)
	}

	co, err := database.GetConundrum(id)
	if err != nil {
		return err
	}

	cs, err := database.GetConundrumComments(id)
	if err != nil {
		return err
	}

	res := &ConundrumResult{
		Conundrum:     co,
		Comments:      cs,
		Next:          co.Id + 1,
		Prev:          co.Id - 1,
		CookiesAgreed: cookiesAgreed(&c),
	}

	return c.Render(http.StatusOK, "conundrum", res)
}

func Star(c echo.Context) error {
	id, err := strconv.Atoi(c.FormValue("id"))
	if err != nil {
		return fmt.Errorf("got malfordmed id: %w", err)
	}

	co, err := database.StarConundrum(id)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "conundrum-stars", co)
}

func UploadResult(c echo.Context) error {
	text := c.FormValue("conundrum")

	err := database.InsertConundrum(text)
	if err != nil {
		return c.Render(http.StatusOK, "upload-form-result", nil)
	}

	return c.Render(http.StatusOK, "upload-form-result", true)
}

func CommentForm(c echo.Context) error {
	com := c.FormValue("comment")
	cid, err := strconv.Atoi(c.FormValue("cid"))
	if err != nil {
		return fmt.Errorf("got malfordmed id: %w", err)
	}

	err = database.InsertComment(cid, com)
	if err != nil {
		return c.Render(http.StatusOK, "comment-form-result", nil)
	}

	return c.Render(http.StatusOK, "comment-form-result", cid)
}

func cookiesAgreed(c *echo.Context) bool {
	_, err := (*c).Cookie("risuhunnik-cookie")
	return err == nil
}
