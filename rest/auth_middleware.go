package rest

import (
	"github.com/labstack/echo/v4"
	"microauth/core"
	"net/http"
)

func NewAuthMiddleware(s core.AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id := core.AuthTokenID(c.Request().Header.Get("Authorization"))
			if err := s.Validate(c.Request().Context(), id); err != nil {
				return c.JSON(
					http.StatusUnauthorized,
					echo.Map{"error": "invalid token"},
				)
			}

			return next(c)
		}
	}
}