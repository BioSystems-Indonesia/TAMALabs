package rest

import (
	"net/http"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	auth_uc "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/auth"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authUC *auth_uc.AuthUseCase
}

func NewAuthHandler(
	authUC *auth_uc.AuthUseCase,
) *AuthHandler {
	return &AuthHandler{
		authUC: authUC,
	}
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req entity.LoginRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	result, err := h.authUC.Login(c.Request().Context(), &req)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, result)
}
