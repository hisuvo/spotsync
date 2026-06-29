package reservations

import (
	"errors"
	"net/http"
	"spotsync/internal/domain/reservations/dto"
	"spotsync/internal/httpresponse"
	"spotsync/internal/middleware"
	"strconv"

	"github.com/labstack/echo/v5"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

// POST /api/v1/reservations
// Access: Authenticated Users (driver, admin)
func (h *handler) Create(c *echo.Context) error {
	claims := middleware.CheckUser(c)

	// Parse userID from JWT string claim
	userID64, err := strconv.ParseUint(claims.UserID, 10, 64)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, httpresponse.ErrorResponse{
			Success: false,
			Message: "Invalid user identity",
			Errors:  map[string]string{"error": err.Error()},
		})
	}

	var req dto.CreateReservationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.ErrorResponse{
			Success: false,
			Message: "Invalid request body",
			Errors:  map[string]string{"error": err.Error()},
		})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.ErrorResponse{
			Success: false,
			Message: "Validation failed",
			Errors:  map[string]string{"error": err.Error()},
		})
	}

	resv, err := h.service.Create(uint(userID64), &req)
	if err != nil {
		if errors.Is(err, ErrZoneFull) {
			return c.JSON(http.StatusConflict, httpresponse.ErrorResponse{
				Success: false,
				Message: "Parking zone is full",
				Errors:  map[string]string{"error": err.Error()},
			})
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.ErrorResponse{
			Success: false,
			Message: "Internal server error",
			Errors:  map[string]string{"error": err.Error()},
		})
	}

	return c.JSON(http.StatusCreated, httpresponse.SuccessResponse{
		Success: true,
		Message: "Reservation confirmed successfully",
		Data:    resv,
	})
}

// GET /api/v1/reservations/my-reservations
// Access: Authenticated Users
func (h *handler) GetMy(c *echo.Context) error {
	claims := middleware.CheckUser(c)

	userID64, err := strconv.ParseUint(claims.UserID, 10, 64)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, httpresponse.ErrorResponse{
			Success: false,
			Message: "Invalid user identity",
			Errors:  map[string]string{"error": err.Error()},
		})
	}

	resvs, err := h.service.GetMyReservations(uint(userID64))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.ErrorResponse{
			Success: false,
			Message: "Internal server error",
			Errors:  map[string]string{"error": err.Error()},
		})
	}

	return c.JSON(http.StatusOK, httpresponse.SuccessResponse{
		Success: true,
		Message: "My reservations retrieved successfully",
		Data:    resvs,
	})
}

// DELETE /api/v1/reservations/:id
// Access: Authenticated Users — drivers can only cancel their own
func (h *handler) Cancel(c *echo.Context) error {
	claims := middleware.CheckUser(c)

	userID64, err := strconv.ParseUint(claims.UserID, 10, 64)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, httpresponse.ErrorResponse{
			Success: false,
			Message: "Invalid user identity",
			Errors:  map[string]string{"error": err.Error()},
		})
	}

	reservationID64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.ErrorResponse{
			Success: false,
			Message: "Invalid reservation ID",
			Errors:  map[string]string{"error": err.Error()},
		})
	}

	// NOTE: In user/service.go, GenerateToken is called as GenerateToken(userID, user.Email, user.Role)
	// so JWT claims.Username = email, and claims.Email = role
	isAdmin := claims.Email == "admin"

	cancelErr := h.service.CancelReservation(uint(userID64), uint(reservationID64), isAdmin)
	if cancelErr != nil {
		if errors.Is(cancelErr, ErrReservationNotFound) {
			return c.JSON(http.StatusNotFound, httpresponse.ErrorResponse{
				Success: false,
				Message: "Reservation not found",
				Errors:  map[string]string{"error": cancelErr.Error()},
			})
		}
		if errors.Is(cancelErr, ErrUnauthorizedCancel) {
			return c.JSON(http.StatusForbidden, httpresponse.ErrorResponse{
				Success: false,
				Message: "Forbidden: you can only cancel your own reservations",
				Errors:  map[string]string{"error": cancelErr.Error()},
			})
		}
		if errors.Is(cancelErr, ErrAlreadyProcessed) {
			return c.JSON(http.StatusBadRequest, httpresponse.ErrorResponse{
				Success: false,
				Message: "Reservation is already completed or cancelled",
				Errors:  map[string]string{"error": cancelErr.Error()},
			})
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.ErrorResponse{
			Success: false,
			Message: "Internal server error",
			Errors:  map[string]string{"error": cancelErr.Error()},
		})
	}

	return c.JSON(http.StatusOK, httpresponse.SuccessResponse{
		Success: true,
		Message: "Reservation cancelled successfully",
	})
}

// GET /api/v1/reservations
// Access: Admin only
func (h *handler) GetAll(c *echo.Context) error {
	claims := middleware.CheckUser(c)
	if claims.Email != "admin" {
		return c.JSON(http.StatusForbidden, httpresponse.ErrorResponse{
			Success: false,
			Message: "Forbidden: admin access required",
			Errors:  map[string]string{"error": "Forbidden"},
		})
	}

	resvs, err := h.service.GetAllReservations()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.ErrorResponse{
			Success: false,
			Message: "Internal server error",
			Errors:  map[string]string{"error": err.Error()},
		})
	}

	return c.JSON(http.StatusOK, httpresponse.SuccessResponse{
		Success: true,
		Message: "All reservations retrieved successfully",
		Data:    resvs,
	})
}