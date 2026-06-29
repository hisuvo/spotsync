package parkingzones

import (
	"spotsync/internal/config"

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

	route := e.Group("/api/v1")

	route.POST("/zones", hdl.Create)
	route.GET("/zones", hdl.GetAll)
	route.GET("/zones/:id", hdl.FindResponseByID)
	route.PUT("/zones/:id", hdl.Update)
	route.DELETE("/zones/:id", hdl.Delete)
}