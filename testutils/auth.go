package testutils

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// CreateAuthMiddleware creates a mock authentication middleware
func CreateAuthMiddleware() func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(c)
		}
	}
}

// Helper function to add a user token to the context
func AddUserToken(c echo.Context, userID int64) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = float64(userID)
	c.Set("user", token)
}
