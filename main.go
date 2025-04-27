package main

import (
	"embed"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strconv"

	"github.com/Zeglius/yafti-go/config"
	"github.com/Zeglius/yafti-go/internal/consts"
	"github.com/Zeglius/yafti-go/ui/pages"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func init() {
	if err := config.LoadConfig(); err != nil {
		log.Panicf("Failed to load config: %v", err)
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
		var screen config.Screen

		sId, err := strconv.Atoi(c.Param("idx"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid screen index")
		}

		if sId < 0 || sId >= len(config.ConfStatus.Screens) {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid screen index")
		}
		screen = config.ConfStatus.Screens[sId]

		handler := newHandler(pages.ActionGroupScreen(screen))
		handler.ServeHTTP(c.Response(), c.Request())

		return nil
	})

	e.GET("/confirm_changes", func(c echo.Context) error {
		// Step 1: Retrieve script IDs from cookie
		scriptIdsCookie, err := c.Cookie("script_ids")
		if err != nil || scriptIdsCookie.Value == "" {
			c.String(http.StatusBadRequest, "Missing cookie 'script_ids'")
			return err
		}

		// Step 2: Parse the JSON from cookie
		var scriptIdsStrs map[string]string
		if err := json.Unmarshal([]byte(scriptIdsCookie.Value), &scriptIdsStrs); err != nil {
			return err
		}

		// Step 3: Extract selected script IDs (those marked as "true")
		scriptIds := make([]string, 0, len(scriptIdsStrs))
		for id, state := range scriptIdsStrs {
			if state == "true" {
				scriptIds = append(scriptIds, id)
			}
		}

		// Step 4: Sort and remove duplicates
		slices.Sort(scriptIds)
		scriptIds = slices.Compact(scriptIds)

		// Step 5: Get actions corresponding to the selected script IDs
		actions, _ := config.ConfStatus.GetActionsByIds(scriptIds)

		// Step 6: Render the confirmation page with selected actions
		handler := newHandler(pages.ConfirmChanges(actions))
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
