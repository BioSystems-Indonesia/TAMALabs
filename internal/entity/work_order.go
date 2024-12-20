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
	ObservationRequests []string        `json:"observation_requests" gorm:"-" validate:"required"`
	PatientIds          []int64         `json:"Patient_ids" gorm:"-" validate:"required"`
	SpecimenIDs         []int64         `json:"Specimen_ids" gorm:"-" validate:"-"`
	CreatedAt           time.Time       `json:"created_at" gorm:"not null"`
	UpdatedAt           time.Time       `json:"updated_at" gorm:"not null"`

	Specimens []Specimen `json:"specimens" gorm:"many2many:work_order_Specimens;->" validate:"-"`
}

type WorkOrderAddSpecimen struct {
	SpecimenIDs []int64 `json:"specimen_ids" gorm:"-" validate:"required"`
}

type WorkOrderGetManyRequest struct {
	GetManyRequest
	SpecimenIDs []int64 `json:"specimen_ids"`
}
