package rest

import (
	"github.com/labstack/echo/v4"
	"microauth/domain"
	"net/http"
	"strconv"
)

type CredentialHandler struct {
	service domain.CredentialService
}

func NewCredentialHandler(s domain.CredentialService) *CredentialHandler {
	return &CredentialHandler{service: s}
}

func (h CredentialHandler) HandleGetCredentials(c echo.Context) error {
	list, err := h.service.ListCredentials(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, list)
}

func (h CredentialHandler) HandleSaveCredential(c echo.Context) error {
	var r domain.SaveParams

	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	if err := h.service.SaveCredential(c.Request().Context(), r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "credential saved"})
}

func (h CredentialHandler) HandleDeleteCredential(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	if err := h.service.RemoveCredential(c.Request().Context(), domain.CredentialID(id)); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "credential deleted"})
}
