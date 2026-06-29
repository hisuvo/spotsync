package server

import (
	"context"
	"net/http"
	"spotsync/internal/config"
	"spotsync/internal/domain/parkingzones"
	"spotsync/internal/domain/user"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"gorm.io/gorm"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func Start(db *gorm.DB, cfg *config.Config) {
	e := echo.New()

	// Register custom validator
	e.Validator = &CustomValidator{validator: validator.New()}

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
		return c.String(http.StatusOK,"Welcome to the SpotSync Backend!")
	})

	user.RegisterRoute(e, db, cfg)
	parkingzones.RegisterRoute(e, db, cfg)

	// Start server
	sc := echo.StartConfig{Address: ":"+ cfg.Port}
	if err := sc.Start(context.Background(), e); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}