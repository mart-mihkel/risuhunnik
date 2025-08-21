package web

import (
	"fmt"
	"net/http"
	"strconv"

	"risuhunnik/pkg/database"

	"github.com/labstack/echo/v4"
)

type ConundrumResult struct {
	Conundrum  *database.Conundrum
	Comments   []database.Comment
	TokenValid bool
	Next       int
	Prev       int
}

func Conundrums(c echo.Context) error {

	conundrums, err := database.GetAllConundrums()
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "conundrums", conundrums)
}

func Conundrum(c echo.Context) error {

	id, err := strconv.Atoi(c.Param("id"))
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

	valid, err := hasValidToken(&c)
	if err != nil {
		return err
	}

	res := &ConundrumResult{
		Conundrum:  conundrum,
		Comments:   comments,
		TokenValid: valid,
		Next:       conundrum.Id + 1,
		Prev:       conundrum.Id - 1,
	}

	return c.Render(http.StatusOK, "conundrum", res)
}
