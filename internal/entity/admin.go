package entity

import (
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type RoleName string

const (
	RoleAdmin   RoleName = "Admin"
	RoleDoctor  RoleName = "Doctor"
	RoleAnalyst RoleName = "Analyst"
)

// Admin represents a user within the system.
type Admin struct {
	ID           int64     `json:"id" gorm:"primaryKey"`
	Fullname     string    `json:"fullname" validate:"required" gorm:"index:idx_admin_fullname"`
	Username     string    `json:"username" validate:"required" gorm:"uniqueIndex:idx_admin_username"`
	Email        *string   `json:"email,omitempty" validate:"omitempty,email" gorm:"uniqueIndex"`
	PasswordHash string    `json:"-"`
	IsActive     bool      `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Password string `json:"password,omitempty" gorm:"-"`

	RolesID []int  `json:"roles_id" gorm:"-"`
	Roles   []Role `json:"roles" gorm:"many2many:admin_roles;"`
}

func (a Admin) GetLastName() string {
	split := strings.Split(a.Fullname, " ")
	if len(split) > 1 {
		return split[len(split)-1]
	}

	return ""
}

func (a Admin) GetFirstName() string {
	split := strings.Split(a.Fullname, " ")
	if len(split) > 1 {
		return strings.Join(split[:len(split)-1], " ")
	}

	return ""
}

func (a Admin) ToAdminClaim(expirationTime time.Time) AdminClaims {
	// Safe handling of Email pointer
	var email string
	if a.Email != nil {
		email = *a.Email
	}

	claims := AdminClaims{
		ID:        a.ID,
		Fullname:  a.Fullname,
		Email:     email,
		IsActive:  a.IsActive,
		Role:      a.Roles[0].Name,
		CreatedAt: a.CreatedAt.Format(time.RFC3339),
		UpdatedAt: a.UpdatedAt.Format(time.RFC3339),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "lims-hl-seven",
		},
	}
	return claims
}

// Role represents a user role within the system.
type Role struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" validate:"required" gorm:"uniqueIndex:idx_role_name_uniq;not null"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Admins []*Admin `json:"admins,omitempty" gorm:"many2many:admin_roles;"`
}

// AdminClaims represents the claims for an admin JWT.
type AdminClaims struct {
	ID        int64  `json:"id"`
	Fullname  string `json:"fullname"`
	Email     string `json:"email"`
	IsActive  bool   `json:"is_active"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	jwt.RegisteredClaims
}

func (ac AdminClaims) ToAdmin() Admin {
	createdAtParse, _ := time.Parse(time.RFC3339, ac.CreatedAt)
	updatedAtParse, _ := time.Parse(time.RFC3339, ac.UpdatedAt)

	email := &ac.Email
	if ac.Email == "" {
		email = nil
	}

	return Admin{
		ID:        ac.ID,
		Fullname:  ac.Fullname,
		Email:     email,
		IsActive:  ac.IsActive,
		CreatedAt: createdAtParse,
		UpdatedAt: updatedAtParse,
	}
}

type GetManyRequestAdmin struct {
	GetManyRequest
	Role []string `query:"role"`
}

type GetManyRequestRole struct {
	GetManyRequest
}
