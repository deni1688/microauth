package rest

import (
	"github.com/labstack/echo/v4"
	"microauth/domain"
	"net/http"
)

type AuthHandler struct {
	service domain.AuthService
}

func NewAuthHandler(s domain.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h AuthHandler) HandleLogin(c echo.Context) error {
	var r domain.AuthParams

	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	t, err := h.service.Authenticate(c.Request().Context(), r)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"token": t})
}

func (h AuthHandler) HandleLogout(c echo.Context) error {
	id := domain.AuthTokenID(c.Request().Header.Get("Authorization"))
	if err := h.service.Expire(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "logout successful"})
}
