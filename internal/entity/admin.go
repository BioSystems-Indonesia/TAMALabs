package entity

import (
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type RoleName string

const (
	RoleAdmin    RoleName = "Admin"
	RoleDoctor   RoleName = "Doctor"
	RoleAnalyzer RoleName = "Analyzer"
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
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	jwt.RegisteredClaims
}

type GetManyRequestAdmin struct {
	GetManyRequest
	Role []string `query:"role"`
}

type GetManyRequestRole struct {
	GetManyRequest
}
