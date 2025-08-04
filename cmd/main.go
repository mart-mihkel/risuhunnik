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

	err = database.ConnectDB("build/risuhunnik.db")
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

	e.Static("/static", "static")

	e.GET("/", web.Index)
	e.GET("/author", web.Author)
	e.GET("/conundrum", web.Conundrum)
	e.GET("/conundrums", web.Conundrums)

	e.POST("/star", web.ToggleStar)
	e.POST("/cookies", web.Cookies)
	e.POST("/upload", web.UploadForm)
	e.POST("/comment", web.CommentForm)

	e.Logger.Fatal(e.Start(":8080"))
}
