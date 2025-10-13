package middleware

import (
	"net/http"

	"github.com/BioSystems-Indonesia/TAMALabs/config"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// simplifiedJWTConfig is a minimal config for JWT middleware.
type simplifiedJWTConfig struct {
	SigningKey []byte
}

// JWTMiddleware struct holds the simplified JWT middleware and its config.
type JWTMiddleware struct {
	config simplifiedJWTConfig
}

// NewJWTMiddleware creates a new JWT middleware instance.
// It takes the signing key and optionally the context key.
func NewJWTMiddleware(cfg *config.Schema) *JWTMiddleware {
	conf := simplifiedJWTConfig{
		SigningKey: []byte(cfg.SigningKey),
	}

	return &JWTMiddleware{config: conf}
}

// Middleware returns the Echo middleware function for JWT validation.
func (m *JWTMiddleware) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("access_token")
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing JWT cookie")
			}

			tokenString := cookie.Value
			if tokenString == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Empty JWT cookie")
			}

			var claims entity.AdminClaims
			token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid JWT signing method")
				}
				return m.config.SigningKey, nil
			})

			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid JWT token")
			}

			if !token.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid JWT token")
			}

			entity.SetEchoContextUser(c, claims)

			return next(c)
		}
	}
}
