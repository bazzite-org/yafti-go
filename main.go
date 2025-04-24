package main

import (
	"bufio"
	"log"
	"net/http"
	"os/exec"

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

	e.POST("/_/command", runCommand)

	// Start server
	log.Printf("Server started at http://localhost:%s", consts.PORT)

	e.Logger.Fatal(e.Start(":" + consts.PORT))
}

func runCommand(c echo.Context) error {
	cmd := c.FormValue("cmd")
	if cmd == "" {
		return c.String(http.StatusBadRequest, "No command provided")
	}

	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().WriteHeader(http.StatusOK)

	commandParts := []string{"/bin/bash", "-c", cmd}
	command := commandParts[0]
	args := commandParts[1:]

	execCmd := exec.Command(command, args...)

	// Get pipes for stdout and stderr
	stdout, err := execCmd.StdoutPipe()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create stdout pipe: "+err.Error())
	}

	stderr, err := execCmd.StderrPipe()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create stderr pipe: "+err.Error())
	}

	if err := execCmd.Start(); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to start command: "+err.Error())
	}

	// Create a channel to signal when command completes
	done := make(chan bool)

	// Stream stdout
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			c.Response().Write([]byte(line + "\n"))
			c.Response().Flush()
		}
	}()

	// Stream stderr
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			c.Response().Write([]byte(line + "\n"))
			c.Response().Flush()
		}
		done <- true
	}()

	// Wait for command to finish
	<-done
	if err := execCmd.Wait(); err != nil {
		c.Response().Write([]byte("data: Command execution failed: " + err.Error() + "\n"))
	}

	c.Response().Flush()

	return nil
}
