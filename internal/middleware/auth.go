package middleware

import (
	"net/http"               // Import Net/HTTP package for response statuses
	"spotsync/internal/auth" // Import Auth package for validating tokens
	"strings"                // Import Strings package to trim string prefixes

	"github.com/labstack/echo/v5" // Import Echo framework
)

// AuthMiddleware is a middleware function that checks for JWT authentication
func AuthMiddleware(JWTService auth.JWTService) echo.MiddlewareFunc{
	// Return the middleware handler closure
	return func(next echo.HandlerFunc) echo.HandlerFunc{
		// Return the context handler closure
		return func(c *echo.Context) error {
			// Retrieve the Authorization header value from the request
			authHeader := c.Request().Header.Get("Authorization")

			// Check if the Authorization header is empty
			if authHeader == "" {
				// Return 401 Unauthorized status with error message if empty
				return c.JSON(http.StatusUnauthorized, "authorization token is missing")
			}

			// Remove the Bearer prefix to get the raw token string
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			// Check if the token string is empty after trimming prefix
			if tokenString == "" {
				// Return 401 Unauthorized status with error message if empty
				return c.JSON(http.StatusUnauthorized, "authorization token is missing")
			}

			// Validate the token and return the claims or error
			claims, err := JWTService.ValidateToken(tokenString)
			// Check if there was an error validating the token
			if err != nil {
				// Return 401 Unauthorized status with error message if token invalid
				return c.JSON(http.StatusUnauthorized, "authorization token is missing")
			}

			// Set the claims as user context in Echo
			c.Set("user", claims)
			// c.Set("id", claims.UserID)
			// c.Set("name", claims.Username)
			// c.Set("email", claims.Email)
			// Call the next handler in the middleware chain
			return next(c)
		}
	}
}

// CheckUser retrieves user claims from Echo context
func CheckUser(c *echo.Context) *auth.Claims {
	// Return user claims from context with type assertion
	return c.Get("user").(*auth.Claims)
}

func CheckMiddleware(name string) echo.MiddlewareFunc{
	return func(next echo.HandlerFunc) echo.HandlerFunc{
		return func(c *echo.Context) error {
			return next(c)
		}
	}
}