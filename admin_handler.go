package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type AdminHandler struct {
	service AdminService
}

func NewAdminHandler(s AdminService) *AdminHandler {
	return &AdminHandler{service: s}
}

func (h AdminHandler) HandleGetAdmins(c echo.Context) error {
	list, err := h.service.ListAdmins(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, list)
}

func (h AdminHandler) HandleSaveAdmin(c echo.Context) error {
	var r SaveParams

	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	if err := h.service.SaveAdmin(c.Request().Context(), r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "admin saved"})
}

func (h AdminHandler) HandleDeleteAdmin(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	if err := h.service.RemoveAdmin(c.Request().Context(), AdminID(id)); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "admin deleted"})
}
