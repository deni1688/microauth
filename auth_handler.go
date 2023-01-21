package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type authHandler struct {
	service AuthService
}

func NewAuthHandler(s AuthService) *authHandler {
	return &authHandler{service: s}
}

func (h authHandler) HandleLogin(c echo.Context) error {
	var r AuthRequest

	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	t, err := h.service.Authenticate(r)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"token": t})
}

func (h authHandler) HandleLogout(c echo.Context) error {
	id := AuthTokenID(c.Request().Header.Get("Authorization"))
	if err := h.service.Invalidate(id); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "logout success"})
}
