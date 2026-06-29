package parkingzones

import (
	"net/http"
	"spotsync/internal/domain/parkingzones/dto"
	"spotsync/internal/httpresponse"
	"spotsync/internal/middleware"
	"strconv"

	"github.com/labstack/echo/v5"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler{
	return &handler{
		service: service,
	}
}

// post /zones
func (h *handler) Create(c *echo.Context) error{
	claims := middleware.CheckUser(c)
	if claims.Email != "admin" {
		return c.JSON(http.StatusForbidden, httpresponse.ErrorResponse{
			Success: false,
			Message: "Forbidden: admin access required",
			Errors:  map[string]string{"error": "Forbidden"},
		})
	}

	var req dto.CreateParkingZoneRequest

	if err := c.Bind(&req); err != nil{
		return c.JSON(http.StatusBadRequest, httpresponse.ErrorResponse{
				Success: false,
				Message: "Invalid request body",
				Errors: map[string]string{"error": err.Error()},
			})
	}

	if err := c.Validate(req); err != nil{
		return c.JSON(http.StatusBadRequest, httpresponse.ErrorResponse{
				Success: false,
				Message: "Validation failed",
				Errors: map[string]string{"error": err.Error()},
			})
	}

	parkingZone, err := h.service.Create(&req)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.ErrorResponse{
				Success: false,
				Message: "Internal server error",
				Errors: map[string]string{"error": err.Error()},
			})
	}

	return c.JSON(http.StatusCreated, httpresponse.SuccessResponse{
		Success: true,
		Message: "Parking zone created successfully",
		Data: parkingZone,
	})
}

func (h *handler) GetAll(c *echo.Context)error{
	parkingZones, err := h.service.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.ErrorResponse{
			Success: false,
			Message: "Internal server error",
			Errors: map[string]string{"error": err.Error()},
		})
	}

	return c.JSON(http.StatusOK, httpresponse.SuccessResponse{
		Success: true,
		Message: "Parking zones retrieved successfully",
		Data: parkingZones,
	})
}

func (h *handler) FindResponseByID(c *echo.Context) error{
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.ErrorResponse{
			Success: false,
			Message: "Invalid parking zone ID",
			Errors: map[string]string{"error": err.Error()},
		})
	}

	parkingZone, err := h.service.FindResponseByID(id)
	if err == ErrParkingZoneNotFound {
		return c.JSON(http.StatusNotFound, httpresponse.ErrorResponse{
			Success: false,
			Message: "Parking zone not found",
			Errors: map[string]string{"error": err.Error()},
		})
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.ErrorResponse{
			Success: false,
			Message: "Internal server error",
			Errors: map[string]string{"error": err.Error()},
		})
	}

	return c.JSON(http.StatusOK, httpresponse.SuccessResponse{
		Success: true,
		Message: "Parking zone retrieved successfully",
		Data: parkingZone,
	})
}

func (h *handler) Update(c *echo.Context) error{
	claims := middleware.CheckUser(c)
	if claims.Email != "admin" {
		return c.JSON(http.StatusForbidden, httpresponse.ErrorResponse{
			Success: false,
			Message: "Forbidden: admin access required",
			Errors:  map[string]string{"error": "Forbidden"},
		})
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.ErrorResponse{
			Success: false,
			Message: "Invalid parking zone ID",
			Errors: map[string]string{"error": err.Error()},
		})
	}

	var req dto.UpdateParkingZoneRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.ErrorResponse{
			Success: false,
			Message: "Invalid request body",
			Errors: map[string]string{"error": err.Error()},
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.ErrorResponse{
			Success: false,
			Message: "Validation failed",
			Errors: map[string]string{"error": err.Error()},
		})
	}

	parkingZone, err := h.service.Update(id, &req)

	if err == ErrParkingZoneNotFound {
		return c.JSON(http.StatusNotFound, httpresponse.ErrorResponse{
			Success: false,
			Message: "Parking zone not found",
			Errors: map[string]string{"error": err.Error()},
		})
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.ErrorResponse{
			Success: false,
			Message: "Internal server error",
			Errors: map[string]string{"error": err.Error()},
		})
	}

	return c.JSON(http.StatusOK, httpresponse.SuccessResponse{
		Success: true,
		Message: "Parking zone updated successfully",
		Data: parkingZone,
	})
}

func (h *handler) Delete(c *echo.Context) error {
	claims := middleware.CheckUser(c)
	if claims.Email != "admin" {
		return c.JSON(http.StatusForbidden, httpresponse.ErrorResponse{
			Success: false,
			Message: "Forbidden: admin access required",
			Errors:  map[string]string{"error": "Forbidden"},
		})
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.ErrorResponse{
			Success: false,
			Message: "Invalid parking zone ID",
			Errors:  map[string]string{"error": err.Error()},
		})
	}

	err = h.service.Delete(id)
	if err == ErrParkingZoneNotFound {
		return c.JSON(http.StatusNotFound, httpresponse.ErrorResponse{
			Success: false,
			Message: "Parking zone not found",
			Errors:  map[string]string{"error": err.Error()},
		})
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.ErrorResponse{
			Success: false,
			Message: "Internal server error",
			Errors:  map[string]string{"error": err.Error()},
		})
	}

	return c.JSON(http.StatusOK, httpresponse.SuccessResponse{
		Success: true,
		Message: "Parking zone deleted successfully",
	})
}