package auth_uc

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/config"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	adminrepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/admin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	adminRepo *adminrepo.AdminRepository
	cfg       *config.Schema

	tokenExpirationTime time.Duration
}

func NewAuthUseCase(adminRepo *adminrepo.AdminRepository, cfg *config.Schema) *AuthUseCase {
	if adminRepo == nil {
		panic("adminRepo cannot be nil")
	}
	if cfg == nil {
		panic("config cannot be nil")
	}

	return &AuthUseCase{
		adminRepo: adminRepo,
		cfg:       cfg,

		tokenExpirationTime: 30 * 24 * time.Hour,
	}
}

// Login function handles admin login logic.
func (u *AuthUseCase) Login(ctx context.Context, req *entity.LoginRequest) (entity.LoginResponse, error) {
	if req.Username == "" {
		return entity.LoginResponse{}, entity.NewHTTPError(http.StatusBadRequest, "Email cannot be empty")
	}
	if req.Password == "" {
		return entity.LoginResponse{}, entity.NewHTTPError(http.StatusBadRequest, "Password cannot be empty")
	}

	admin, err := u.adminRepo.FindOneByUsername(ctx, req.Username)
	if err != nil {
		return entity.LoginResponse{}, entity.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	if !admin.IsActive {
		return entity.LoginResponse{}, entity.NewHTTPError(http.StatusUnauthorized, "Account is not active")
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(req.Password))
	if err != nil {
		return entity.LoginResponse{}, entity.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	token, err := u.createAccessToken(ctx, admin)
	if err != nil {
		return entity.LoginResponse{}, fmt.Errorf("failed to create access token: %w", err)
	}

	return entity.LoginResponse{
		Admin:       admin,
		AccessToken: token,
	}, nil
}

// createAccessToken generates a JWT access token for the given admin.
// The token will be expired after 30 days.
func (u *AuthUseCase) createAccessToken(ctx context.Context, admin entity.Admin) (string, error) {
	expirationTime := time.Now().Add(u.tokenExpirationTime)

	claims := admin.ToAdminClaim(expirationTime)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Check if signing key is available
	if u.cfg.SigningKey == "" {
		return "", fmt.Errorf("signing key is not configured")
	}

	signedToken, err := token.SignedString([]byte(u.cfg.SigningKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT token: %w", err)
	}

	return signedToken, nil
}
