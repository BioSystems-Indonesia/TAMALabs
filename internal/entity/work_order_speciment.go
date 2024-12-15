package entity

import (
	"time"
)

type WorkOrderSpeciment struct {
	WorkOrderID int64     `json:"work_order_id" gorm:"not null" validate:"required"`
	SpecimentID int64     `json:"speciment_id" gorm:"not null" validate:"required"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"not null"`
}
