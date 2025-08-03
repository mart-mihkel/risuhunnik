package web

import (
	"fmt"
	"net/http"
	"strconv"

	"risuhunnik/pkg/database"

	"github.com/labstack/echo/v4"
)

// TODO: look into wrapping cookies agreed in a middleware/general result

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

	c.SetCookie(&http.Cookie{Name: "agreed"})

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

	res := &ConundrumResult{
		Conundrum:     conundrum,
		Comments:      comments,
		Next:          conundrum.Id + 1,
		Prev:          conundrum.Id - 1,
		CookiesAgreed: cookiesAgreed(&c),
	}

	return c.Render(http.StatusOK, "conundrum", res)
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
	comment := c.FormValue("comment")
	id, err := strconv.Atoi(c.FormValue("conundrum-id"))
	if err != nil {
		return fmt.Errorf("got malfordmed id: %w", err)
	}

	err = database.InsertComment(id, comment)
	if err != nil {
		return c.Render(http.StatusOK, "comment-form-result", nil)
	}

	return c.Render(http.StatusOK, "comment-form-result", id)
}

func cookiesAgreed(c *echo.Context) bool {
	_, err := (*c).Cookie("agreed")
	return err == nil
}
