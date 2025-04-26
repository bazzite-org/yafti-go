package main

import (
	"embed"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Zeglius/yafti-go/config"
	"github.com/Zeglius/yafti-go/internal/consts"
	"github.com/Zeglius/yafti-go/ui/pages"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var conf *config.Config

func init() {
	if c, err := config.LoadConfig(); err != nil {
		panic(err)
	} else {
		conf = c
	}
}

// Default [templ.Handler] with streaming enabled by default
func newHandler(c templ.Component, options ...func(*templ.ComponentHandler)) *templ.ComponentHandler {
	opts := []func(*templ.ComponentHandler){templ.WithStreaming()}
	opts = append(opts, options...)
	return templ.Handler(c, opts...)
}

//go:embed static/**
var static embed.FS

func runServer() error {
	e := echo.New()

	e.Use(middleware.Logger())

	// Set up static file serving
	fs := echo.MustSubFS(static, "static")
	e.StaticFS("/static/", fs)

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

	e.GET("/action_group/:idx", func(c echo.Context) error {
		var screen *config.Screen

		sId, err := strconv.Atoi(c.Param("idx"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid screen index")
		}

		if sId < 0 || sId >= len(conf.Screens) {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid screen index")
		}
		screen = &conf.Screens[sId]

		handler := newHandler(pages.ActionGroupScreen(screen.Actions))
		handler.ServeHTTP(c.Response(), c.Request())

		return nil
	})

	e.POST("/_/apply_changes", func(c echo.Context) error {
		// TODO: Use "{script_ids: [<id>]}" to obtain script ids, that we
		// use to extract scripts contents from a config struct and pass them
		// as a `[]string` to [pages.ApplyChanges].

		type Payload struct {
			ScriptIds []string `form:"script_ids"`
		}

		payload := Payload{}
		if err := c.Bind(&payload); err != nil {
			return err
		}

		// TODO: Replace placeholders with actual values using the payload data

		cmds := []string{
			`echo "placeholder 1"`,
			`echo "placeholder 2"`,
			`echo "placeholder 3"`,
		}

		handler := newHandler(pages.ApplyChanges(cmds))

		handler.ServeHTTP(c.Response(), c.Request())

		return nil
	})

	e.POST("/_/post_test", func(c echo.Context) error {
		data := struct {
			POSTParams url.Values        `json:"POST_params"`
			Cookies    map[string]string `json:"cookies"`
		}{}

		if v, err := c.FormParams(); err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		} else {
			data.POSTParams = v
		}

		if c.Cookies() != nil {
			data.Cookies = make(map[string]string)
			for _, v := range c.Cookies() {
				k := v.Name
				data.Cookies[k] = v.Value
			}
		}

		return c.JSON(http.StatusOK, data)
	})

	// Start server
	log.Printf("Server started at http://localhost:%s", consts.PORT)

	return e.Start(":" + consts.PORT)
}

func main() {
	if err := runServer(); err != nil {
		log.Fatal(err)
	}
}
