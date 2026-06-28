package user

import (
	"spotsync/internal/auth"
	"spotsync/internal/config"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoute(e *echo.Echo, db *gorm.DB, cfg *config.Config) {
	repo := NewRepository(db)
	jwt := auth.NewJWTService(cfg.JWTSecret, cfg.TokenDuration)
	svc := NewService(repo, jwt)
	hdl := NewHandler(svc, cfg)

	route := e.Group("/api/v1/auth")

	route.POST("/register", hdl.CreateUser)
	route.POST("/login", hdl.LoginUser)
}