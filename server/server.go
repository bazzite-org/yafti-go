package server

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/Zeglius/yafti-go/config"
	"github.com/Zeglius/yafti-go/internal/consts"
	"github.com/Zeglius/yafti-go/ui/pages"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	e            *echo.Echo
	lastBeat     time.Time
	m            sync.Mutex
	shutdownCtx  context.Context
	cancel       context.CancelFunc
	StaticAssets *embed.FS // This var is set in main.go
}

// Default [templ.Handler] with streaming enabled by default
func newHandler(c templ.Component, options ...func(*templ.ComponentHandler)) *templ.ComponentHandler {
	opts := []func(*templ.ComponentHandler){templ.WithStreaming()}
	opts = append(opts, options...)
	return templ.Handler(c, opts...)
}

func New() *Server {
	e := echo.New()
	ctx, cancel := context.WithCancel(context.Background())

	return &Server{
		e:           e,
		lastBeat:    time.Now(),
		shutdownCtx: ctx,
		cancel:      cancel,
	}
}

func (s *Server) heartbeatHandler(c echo.Context) error {
	s.m.Lock()
	s.lastBeat = time.Now()
	s.m.Unlock()
	return c.String(http.StatusOK, "Hearthbeat received")
}

func (s *Server) monitorHeartbeat() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.shutdownCtx.Done():
			return
		case <-ticker.C:
			s.m.Lock()
			if time.Since(s.lastBeat) > consts.HEARTBEAT_SECONDS*time.Second {
				log.Printf("No heartbeat for %d seconds, shutting down server\n", consts.HEARTBEAT_SECONDS)
				s.m.Unlock()
				s.cancel()
				if err := s.e.Shutdown(context.Background()); err != nil {
					log.Printf("Shutdown error: %v", err)
				}
				return
			}
			s.m.Unlock()
		}
	}
}

func (s *Server) Start() error {
	e := s.e

	e.Use(middleware.Logger())

	// Set up static file serving
	if s.StaticAssets == nil {
		return errors.New("StaticAssets is not populated. Ensure it is set in main.go")
	}
	fs := echo.MustSubFS(*s.StaticAssets, "static")
	e.StaticFS("/static/", fs)

	// Handle heartbeat, so we shutdown the server automatically
	// when there is no client connected over a period of time.
	e.GET("/_/heartbeat", s.heartbeatHandler)

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

	// Update POST handler for confirm_changes
	e.POST("/confirm_changes", func(c echo.Context) error {
		var scriptIdsStrs map[string]string

		// First try to get data from form submission
		formValue := c.FormValue("scriptIds")
		if formValue != "" {
			if err := json.Unmarshal([]byte(formValue), &scriptIdsStrs); err == nil {
				// Successfully got data from form submission
				log.Println("Using script IDs from form data")
			}
		}

		// If that didn't work, fall back to cookie
		if scriptIdsStrs == nil {
			scriptIdsCookie, cookieErr := c.Cookie("script_ids")
			if cookieErr == nil && scriptIdsCookie.Value != "" {
				if err := json.Unmarshal([]byte(scriptIdsCookie.Value), &scriptIdsStrs); err == nil {
					// Successfully got data from cookie
					log.Println("Using script IDs from cookie")
				} else {
					log.Printf("Error parsing cookie value: %v", err)
				}
			} else {
				log.Printf("Cookie error or empty: %v", cookieErr)
			}
		}

		// If we still don't have data, return error
		if scriptIdsStrs == nil {
			return c.String(http.StatusBadRequest, "No script IDs found in form data or cookies")
		}

		// Continue with the existing logic
		scriptIds := make([]string, 0, len(scriptIdsStrs))
		for id, state := range scriptIdsStrs {
			if state == "true" {
				scriptIds = append(scriptIds, id)
			}
		}

		slices.Sort(scriptIds)
		scriptIds = slices.Compact(scriptIds)

		actions, _ := config.ConfStatus.GetActionsByIds(scriptIds)

		handler := newHandler(pages.ConfirmChanges(actions))
		handler.ServeHTTP(c.Response(), c.Request())

		return nil
	})

	e.POST("/_/apply_changes", func(c echo.Context) error {
		// Get script IDs from the request payload
		type Payload struct {
			ScriptIds []string `form:"script_ids"`
		}

		payload := Payload{}
		if err := c.Bind(&payload); err != nil {
			log.Printf("Failed to bind payload: %v", err)
			return c.String(http.StatusBadRequest, "Invalid request format")
		}

		if len(payload.ScriptIds) == 0 {
			log.Printf("No script IDs provided in request")
			return c.String(http.StatusBadRequest, "No script IDs provided")
		}

		// Get actions corresponding to the script IDs
		actions, found := config.ConfStatus.GetActionsByIds(payload.ScriptIds)
		if !found || len(actions) == 0 {
			log.Printf("No actions found for the provided script IDs")
			return c.String(http.StatusBadRequest, "No actions found for the provided script IDs")
		}

		// Extract script commands from the actions
		cmds := make([]string, 0, len(actions))
		for _, action := range actions {
			if action.Script != "" {
				cmds = append(cmds, action.Script)
			}
		}

		if len(cmds) == 0 {
			log.Printf("No scripts found in the selected actions")
			return c.String(http.StatusBadRequest, "Selected actions contain no scripts to execute")
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
	go s.monitorHeartbeat()
	log.Printf("Server started at http://localhost:%s", consts.PORT)
	return s.e.Start("127.0.0.1:" + consts.PORT)
}
