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
	ID                  int64           `json:"id,omitempty" gorm:"primaryKey;autoIncrement"`
	Status              WorkOrderStatus `json:"status" gorm:"not null" validate:"work-order-status"`
	ObservationRequests []string        `json:"observation_requests,omitempty" gorm:"-" validate:"required"`
	PatientIDs          []int64         `json:"patient_ids,omitempty" gorm:"-" validate:"required"`
	DeviceID            int64           `json:"device_id" gorm:"not null"`
	CreatedAt           time.Time       `json:"created_at" gorm:"not null"`
	UpdatedAt           time.Time       `json:"updated_at" gorm:"not null"`

	Patient  []Patient  `json:"patient_list,omitempty" gorm:"many2many:work_order_patients;->" validate:"-"`
	Specimen []Specimen `json:"specimen_list,omitempty" gorm:"foreignKey:OrderID;->" validate:"-"`
	Device   Device     `json:"device" gorm:"->" validate:"-"`
}

type WorkOrderRunRequest struct {
	DeviceID    int64 `json:"device_id" gorm:"-" validate:"required"`
	WorkOrderID int64 `json:"work_order_id" gorm:"-" validate:"required"`
}

type WorkOrderCancelRequest struct {
	WorkOrderID int64 `json:"work_order_id" gorm:"-" validate:"required"`
}

type WorkOrderGetManyRequest struct {
	GetManyRequest
	SpecimenIDs []int64 `json:"specimen_ids"`
}
