package entity

import "time"

type WorkOrderStatus string

const (
	WorkOrderStatusNew       WorkOrderStatus = "NEW"
	WorkOrderStatusPending   WorkOrderStatus = "PENDING"
	WorkOrderCancelled       WorkOrderStatus = "CANCELLED"
	WorkOrderStatusCompleted WorkOrderStatus = "SUCCESS"
)

type WorkOrder struct {
	ID           int64           `json:"id,omitempty" gorm:"primaryKey;autoIncrement"`
	Description  string          `json:"description" gorm:"not null;default:''" validate:"required"`
	Status       WorkOrderStatus `json:"status" gorm:"not null" validate:"work-order-status"`
	SpecimentIDs []int64         `json:"speciment_ids" gorm:"-" validate:"required"`
	CreatedAt    time.Time       `json:"created_at" gorm:"not null"`
	UpdatedAt    time.Time       `json:"updated_at" gorm:"not null"`

	Speciments []Specimen `json:"speciments" gorm:"many2many:work_order_speciments;->" validate:"-"`
}

type WorkOrderAddSpeciment struct {
	SpecimentIDs []int64 `json:"speciment_ids" gorm:"-" validate:"required"`
}

type WorkOrderGetManyRequest struct {
	GetManyRequest
	SpecimentIDs []int64 `json:"speciment_ids"`
}
