package main

import (
	"log"
	"net/http"

	"github.com/Zeglius/yafti-go/internal/consts"
	"github.com/Zeglius/yafti-go/ui/pages"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Default [templ.Handler] with streaming enabled by default
func newHandler(c templ.Component, options ...func(*templ.ComponentHandler)) *templ.ComponentHandler {
	opts := []func(*templ.ComponentHandler){templ.WithStreaming()}
	opts = append(opts, options...)
	return templ.Handler(c, opts...)
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())

	// Set up static file serving
	e.Static("/static/", "./static")

	// Handle pages routes
	e.GET("/", echo.WrapHandler(
		newHandler(pages.Home()),
	))

	e.GET("/about", func(c echo.Context) error {
		return c.NoContent(http.StatusNotFound)
	})

	e.GET("/_/dummy", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct {
			Name    string   `json:"name"`
			Age     int      `json:"age"`
			Hobbies []string `json:"hobbies"`
		}{
			Name:    "John Doe",
			Age:     30,
			Hobbies: []string{"Reading", "Hiking", "Cooking"},
		})
	})

	e.POST("/_/apply_changes", func(c echo.Context) error {
		cmds := []string{
			`echo "placeholder 1"`,
			`echo "placeholder 2"`,
			`echo "placeholder 3"`,
		}

		handler := newHandler(pages.ApplyChanges(cmds))

		handler.ServeHTTP(c.Response(), c.Request())

		return nil
	})

	// Start server
	log.Printf("Server started at http://localhost:%s", consts.PORT)

	e.Logger.Fatal(e.Start(":" + consts.PORT))
}
