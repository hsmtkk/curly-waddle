package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/hsmtkk/curly-waddle/trans"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	portStr := mustEnvVar("PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("failed to parse PORT env var as int; %s; %v", portStr, err)
	}

	chanAccessToken := mustEnvVar("CHANNEL_ACCESS_TOKEN")

	hdl := newHandler(chanAccessToken, trans.New())

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/", hdl.Handle)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}

func mustEnvVar(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("you must define env var: %s", key)
	}
	return val
}
