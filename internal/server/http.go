package server

import (
	"context"
	"net/http"
	"spotsync/internal/config"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"gorm.io/gorm"
)

func Start(db *gorm.DB, env *config.Config) {
	e := echo.New()

	// Middleware
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	// CORS default
	// Allows request from any origin with GET, POST, PUT, HEAD OR DELETE method
	e.Use(middleware.CORS("*"))

	// CORS restricted
	// Allows requests from any `https://spotsync.com` or `https://spotsync.net` origin
	// e.Use(middleware.CORS("https://spotsync.com", "https://spotsync.net"))

	// Routees
	e.GET("/", func(c *echo.Context) error {
		
		return c.JSON(http.StatusOK,"Welcome to the SpotSync Backend!")
	})

	// Start server
	sc := echo.StartConfig{Address: ":"+ env.Port}
	if err := sc.Start(context.Background(), e); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}