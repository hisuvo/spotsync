package user

import (
	"net/http"
	"spotsync/internal/config"
	"spotsync/internal/domain/user/dto"
	"spotsync/internal/httpresponse"

	"github.com/labstack/echo/v5"
)

type handler struct {
	service UserService
	cfg     *config.Config
}

func NewHandler(service UserService, cfg *config.Config) *handler {
	return &handler{
		service: service,
		cfg:     cfg,
	}
}

func (h *handler) CreateUser(c *echo.Context) error {
	var req dto.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	res, err := h.service.CreateUser(&req)

	if err != nil {
		if err == ErrorAlreadyExist {
			return c.JSON(http.StatusConflict, httpresponse.ErrorResponse{
				Success: false,
				Code: http.StatusConflict,
				Message: "User already exists",
				Errors: map[string]string{"error": err.Error()},
			})
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.ErrorResponse{
				Success: false,
				Code: http.StatusInternalServerError,
				Message: "Internal server error",
				Errors: map[string]string{"error": err.Error()},
			})
	}

	return c.JSON(http.StatusCreated, httpresponse.SuccessResponse{
		Success: true,
		Code: http.StatusCreated,
		Message: "User registered successfully",
		Data: res,
	})
}

func (h *handler) LoginUser(c *echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	res, err := h.service.LoginUser(req.Email, req.Password)

	if err != nil {
		if err == ErrInvalideCredentials {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, httpresponse.SuccessResponse{
		Success: true,
		Code: http.StatusOK,
		Message: "Login successful",
		Data: res,
	})
}