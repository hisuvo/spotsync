package parkingzones

import (
	"spotsync/internal/auth"
	"spotsync/internal/config"
	"spotsync/internal/middleware"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// RegisterRoute wires up the parkingzones domain and registers all routes.
// reservationCounter is provided by the caller (server layer) to avoid an import cycle:
// reservations → parkingzones, so parkingzones must never import reservations.
func RegisterRoute(e *echo.Echo, db *gorm.DB, cfg *config.Config, reservationCounter ReservationCounter) {

	repo := NewParkingZoneRepository(db)
	svc := NewParkingZoneService(repo, reservationCounter)
	hdl := NewHandler(svc)

	jwtService := auth.NewJWTService(cfg.JWTSecret, cfg.TokenDuration)
	authMiddleware := middleware.AuthMiddleware(jwtService)

	route := e.Group("/api/v1")

	// Protected routes (admin/authenticated check inside handler)
	authRoute := e.Group("/api/v1", authMiddleware)
	authRoute.POST("/zones", hdl.Create)
	authRoute.PUT("/zones/:id", hdl.Update)
	authRoute.DELETE("/zones/:id", hdl.Delete)

	// Public routes
	route.GET("/zones", hdl.GetAll)
	route.GET("/zones/:id", hdl.FindResponseByID)
}