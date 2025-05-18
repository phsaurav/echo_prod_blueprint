package server

import (
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	errs "github.com/phsaurav/echo_prod_blueprint/pkg/error"
	"github.com/phsaurav/echo_prod_blueprint/pkg/response"
)

// JWTAuth middleware validates JWT tokens and adds user info to context
func JWTAuth(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get token from Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return response.ErrorBuilder(errs.Unauthorized(errors.New("missing authorization header"))).Send(c)
			}

			var tokenString string

			// Check if the Authorization header has the format "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				// Extract token from "Bearer <token>" format
				tokenString = parts[1]
			} else {
				// Use the entire header as the token
				tokenString = authHeader
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(secret), nil
			})

			if err != nil {
				return response.ErrorBuilder(errs.Unauthorized(err)).Send(c)
			}

			// Check if token is valid
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				c.Set("user", token)
				c.Set("user_id", int64(claims["user_id"].(float64)))
				return next(c)
			}

			return response.ErrorBuilder(errs.Unauthorized(errors.New("invalid token"))).Send(c)
		}
	}
}
