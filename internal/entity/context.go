package entity

import (
	"github.com/labstack/echo/v4"
)

const (
	// contextKeyUser is a key for user in context
	contextKeyUser = "user"
)

func SetEchoContextUser(c echo.Context, claims AdminClaims) {
	c.Set(contextKeyUser, claims)
}

func GetEchoContextUser(c echo.Context) AdminClaims {
	user := c.Get(contextKeyUser)
	if user == nil {
		return AdminClaims{}
	}

	claims, ok := user.(AdminClaims)
	if !ok {
		return AdminClaims{}
	}

	return claims
}
