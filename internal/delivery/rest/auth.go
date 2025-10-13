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

	cookie := &http.Cookie{
		Name:     "access_token",
		Value:    result.AccessToken,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
	}

	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, result)
}

func (h *AuthHandler) GetProfile(c echo.Context) error {
	claims := entity.GetEchoContextUser(c)
	if claims.ID == 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user context")
	}

	admin := claims.ToAdmin()
	return c.JSON(http.StatusOK, map[string]interface{}{
		"admin": admin,
		"user":  admin,
	})
}

func (h *AuthHandler) GetPermissions(c echo.Context) error {
	claims := entity.GetEchoContextUser(c)
	if claims.ID == 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user context")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"role":       claims.Role,
		"permission": claims.Role,
	})
}

func (h *AuthHandler) Logout(c echo.Context) error {
	// Clear the access_token cookie by setting it to empty and expired
	cookie := &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
		MaxAge:   -1, // This will delete the cookie
	}

	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})
}
