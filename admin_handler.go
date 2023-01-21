package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type adminHandler struct {
	service AdminService
}

func NewAdminHandler(s AdminService) *adminHandler {
	return &adminHandler{service: s}
}

func (h adminHandler) HandleGetAdmins(c echo.Context) error {
	list, err := h.service.ListAdmins()
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, list)
}

func (h adminHandler) HandleSaveAdmin(c echo.Context) error {
	var r SaveRequest

	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	if err := h.service.SaveAdmin(r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "admin saved"})
}

func (h adminHandler) HandleDeleteAdmin(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	if err := h.service.RemoveAdmin(AdminID(id)); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "admin deleted"})
}
