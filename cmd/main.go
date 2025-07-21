package main

import (
	"html/template"
	"io"
	"log"

	"risuhunnik/pkg/database"
	"risuhunnik/pkg/pages"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data any, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	tmpls, err := template.New("").ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("couldn't initialize templates: %v", err)
	}

	err = database.ConnectDB("risuhunnik.db")
	if err != nil {
		log.Fatalf("couldn't initialize db: %v", err)
	}

	e := echo.New()
	e.Renderer = &TemplateRenderer{templates: tmpls}

	e.Use(middleware.Logger())
	e.Static("/css", "css")

	e.GET("/", pages.Index)
	e.GET("/tags", pages.Tags)
	e.GET("/modal", pages.Modal)
	e.GET("/conundrums", pages.Conundrums)

	e.POST("/search", pages.SearchConundrums)

	e.Logger.Fatal(e.Start(":8080"))
}
