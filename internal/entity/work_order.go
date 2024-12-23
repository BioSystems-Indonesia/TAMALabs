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
	ID         int64           `json:"id,omitempty" gorm:"primaryKey;autoIncrement"`
	Status     WorkOrderStatus `json:"status" gorm:"not null" validate:"work-order-status"`
	PatientIDs []int64         `json:"patient_ids" gorm:"-" validate:"required"`
	CreatedAt  time.Time       `json:"created_at" gorm:"not null"`
	UpdatedAt  time.Time       `json:"updated_at" gorm:"not null"`

	Patient  []Patient  `json:"patient_list" gorm:"many2many:work_order_patients;->" validate:"-"`
	Specimen []Specimen `json:"specimen_list" gorm:"foreignKey:OrderID;->" validate:"-"`
}

type WorkOrderRunRequest struct {
	WorkOrderID []int64 `json:"work_order_ids" gorm:"-" validate:"required"`
}

type WorkOrderGetManyRequest struct {
	GetManyRequest
	SpecimenIDs []int64 `json:"specimen_ids"`
}
