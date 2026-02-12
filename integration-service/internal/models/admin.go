package models

import "time"

type Admin struct {
	ID           uint
	Fullname     string
	Username     string
	Email        *string
	PasswordHash string `gorm:"column:password_hash"`
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time

	Roles []Role `gorm:"many2many:admin_roles"`
}

type Role struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" validate:"required" gorm:"uniqueIndex:idx_role_name_uniq;not null"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Admins []*Admin `json:"admins,omitempty" gorm:"many2many:admin_roles;"`
}

func (Admin) TableName() string {
	return "admins"
}
