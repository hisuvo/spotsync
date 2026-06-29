package reservations

import (
	"spotsync/internal/auth"
	"spotsync/internal/config"
	"spotsync/internal/middleware"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// RegisterRoute wires up the reservations domain and registers all routes.
func RegisterRoute(e *echo.Echo, db *gorm.DB, cfg *config.Config) {
	repo := NewRepository(db)
	svc := NewService(repo)
	hdl := NewHandler(svc)

	jwtService := auth.NewJWTService(cfg.JWTSecret, cfg.TokenDuration)
	authMiddleware := middleware.AuthMiddleware(jwtService)

	// Routes that require any authenticated user (driver or admin)
	authGroup := e.Group("/api/v1/reservations", authMiddleware)
	
	authGroup.POST("", hdl.Create)
	authGroup.GET("/my-reservations", hdl.GetMy)
	authGroup.DELETE("/:id", hdl.Cancel)

	// Admin-only route: get ALL reservations in the system
	// Uses the same authMiddleware — role enforcement is handled inside the handler
	authGroup.GET("", hdl.GetAll)
}