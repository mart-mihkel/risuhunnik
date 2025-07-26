package main

import (
	"html/template"
	"io"
	"log"

	"risuhunnik/pkg/database"
	"risuhunnik/pkg/web"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
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

	ms := middleware.NewRateLimiterMemoryStore(rate.Limit(20))
	gz := middleware.GzipConfig{Level: 5}

	e.Use(middleware.GzipWithConfig(gz))
	e.Use(middleware.RateLimiter(ms))
	e.Use(middleware.Logger())

	e.Static("/css", "css")

	e.GET("/", web.Index)
	e.GET("/tags", web.Tags)
	e.GET("/modal", web.Modal)
	e.GET("/conundrums", web.Conundrums)

	e.POST("/add", web.AddModal)
	e.POST("/star", web.StarButton)
	e.POST("/search", web.SearchModal)

	e.Logger.Fatal(e.Start(":8080"))
}
