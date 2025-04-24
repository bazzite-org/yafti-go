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

func main() {
	e := echo.New()

	e.Use(middleware.Logger())

	// Set up static file serving
	e.Static("/static/", "./static")

	// Handle pages routes
	e.GET("/", echo.WrapHandler(
		templ.Handler(pages.Home()),
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

	e.GET("/_/command", func(c echo.Context) error {
		// Placeholder

		c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
		c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
		c.Response().Header().Set(echo.HeaderConnection, "keep-alive")

		handler := templ.Handler(pages.Command(`for f in *; do echo $f; sleep 1; done`), templ.WithStreaming())
		handler.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	e.POST("/_/command", func(c echo.Context) error {

		c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
		c.Response().Header().Set(echo.HeaderConnection, "keep-alive")

		cmd := c.FormValue("cmd")

		if cmd == "" {
			return c.String(http.StatusBadRequest, "Missing \"cmd\" parameter")
		}

		handler := templ.Handler(pages.Command(cmd), templ.WithStreaming())
		handler.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	// Start server
	log.Printf("Server started at http://localhost:%s", consts.PORT)

	e.Logger.Fatal(e.Start(":" + consts.PORT))
}
