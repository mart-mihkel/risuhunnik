package web

import (
	"fmt"
	"net/http"
	"strconv"

	"risuhunnik/pkg/database"

	"github.com/labstack/echo/v4"
)

type ConundrumResult struct {
	Conundrum *database.Conundrum
	Comments  []database.Comment
	Next      int
	Prev      int
}

func Index(c echo.Context) error {
	cs, err := database.GetAllConundrums()
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "index.html", cs)
}

func Conundrums(c echo.Context) error {
	cs, err := database.GetAllConundrums()
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "conundrums", cs)
}

func Conundrum(c echo.Context) error {
	strid := c.QueryParam("id")
	id, err := strconv.Atoi(strid)
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
		Conundrum: co,
		Comments:  cs,
		Next:      co.Id + 1,
		Prev:      co.Id - 1,
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
	co := c.FormValue("comment")
	strid := c.FormValue("cid")
	cid, err := strconv.Atoi(strid)
	if err != nil {
		return fmt.Errorf("got malfordmed id: %w", err)
	}

	err = database.InsertComment(cid, co)
	if err != nil {
		return c.Render(http.StatusOK, "comment-form-result", nil)
	}

	return c.Render(http.StatusOK, "comment-form-result", cid)
}
