package main

import (
	"io"
	"net/http"
	"text/template"

	. "github.com/Zeglius/yafti-go/internal/consts"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data any, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplateRenderer() *TemplateRenderer {
	tmpls := template.Must(template.ParseGlob(HTML_TMPL_PATH + "/*.html"))
	return &TemplateRenderer{tmpls}
}

type Count struct {
	Count int
}

func main() {

	e := echo.New()
	e.Use(middleware.Logger())
	e.HideBanner = true

	count := Count{0}
	e.Renderer = newTemplateRenderer()

	e.GET("/", func(c echo.Context) error {
		count.Count++
		return c.Render(http.StatusOK, "index", count)
	})

	e.Static("/static", "static")

	e.Logger.Fatal(e.Start(":" + PORT))
}
