package models

import "time"

type ObservationRequest struct {
	ID              int64
	TestCode        string
	TestTypeID      *int
	TestDescription string
	RequestedDate   time.Time
	ResultStatus    string
	SpecimenID      int64
	CreatedAt       time.Time
	UpdatedAt       time.Time

	TestType  TestType  `gorm:"foreignKey:TestTypeID;references:ID;->"`
	WorkOrder WorkOrder `gorm:"-"`
}
