package models

import "time"

type ObservationResult struct {
	ID             int64
	SpecimenID     int64
	Code           string
	Description    string
	Values         []string `gorm:"serializer:json"`
	Type           string
	Unit           string
	ReferenceRange string
	Date           time.Time
	AbnormalFlag   []string `gorm:"serializer:json"`
	Comments       string
	Picked         bool
	CreatedBy      int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
	TestTypeID     *int
}
