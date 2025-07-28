package web

import (
	"net/http"

	"risuhunnik/pkg/database"

	"github.com/labstack/echo/v4"
)

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

func UploadForm(c echo.Context) error {
	return c.Render(http.StatusOK, "upload-form", nil)
}

func UploadResult(c echo.Context) error {
	text := c.FormValue("conundrum")
	co := &database.Conundrum{Text: text}

	_, err := database.InsertConundrum(co)
	if err != nil {
		return c.Render(http.StatusOK, "upload-form-result", "upload failed!")
	}

	return c.Render(http.StatusOK, "upload-form-result", "conundrum uploaded!")
}
