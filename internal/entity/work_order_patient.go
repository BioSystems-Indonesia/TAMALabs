package entity

import (
	"time"
)

type WorkOrderPatient struct {
	WorkOrderID int64     `json:"work_order_id" gorm:"not null;index:work_order_patient_uniq,unique" validate:"required"`
	PatientID   int64     `json:"patient_id" gorm:"not null;index:work_order_patient_uniq,unique" validate:"required"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"not null"`
}
