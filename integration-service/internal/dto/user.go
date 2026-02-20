package dto

import "time"

type UserPushItem struct {
	UserId       uint      `json:"external_id"`
	FullName     string    `json:"fullname"`
	Username     string    `json:"username"`
	Email        *string   `json:"email"`
	PasswordHash string    `json:"password_hash"`
	IsActive     bool      `json:"is_active"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserRequest struct {
	UserID   string `json:"user_id" validate:"required"`
	LabID    string `json:"lab_id" validate:"required"`
	Fullname string `json:"fullname" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
	Role     uint   `json:"role" validate:"required"`
	Status   string `json:"status"`
}
