package parkingzones

import (
	"spotsync/internal/config"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoute(e *echo.Echo, db *gorm.DB, cfg *config.Config){
	
	repo := NewParkingZoneRepository(db)
	svc := NewParkingZoneService(repo)
	hdl := NewHandler(svc)
	
	route := e.Group("/api/v1")

	route.POST("/zones", hdl.Create)
	route.GET("/zones", hdl.GetAll)
	route.GET("/zones/:id", hdl.FindById)
	// route.PUT("/zones/:id", hdl.Update)
	// route.DELETE("/zones/:id", hdl.Delete)
}