package entity

import (
	"time"
)

type Speciment struct {
	ID          int64     `json:"id,omitempty" gorm:"primaryKey;autoIncrement"`
	Description string    `json:"description" gorm:"not null;default:''" validate:"required"`
	Type        string    `json:"type" gorm:"not null" validate:"required,speciment-type"`
	Test        string    `json:"test" gorm:"not null" validate:"required,speciment-test"`
	PatientID   int64     `json:"patient_id" gorm:"not null" validate:"required"`
	Barcode     string    `json:"barcode" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"not null"`

	// Relationships
	Patient Patient `json:"patient" gorm:"foreignKey:PatientID" validate:"-"`
}

type SpecimentGetManyRequest struct {
	GetManyRequest
	PatientID int64 `query:"patient_id"`
}
