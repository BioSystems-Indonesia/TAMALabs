package models

import "time"

type Specimen struct {
	ID             int
	PatientID      int
	OrderID        int
	Type           string
	CollectionDate string
	ReceivedDate   time.Time
	Source         string
	Condition      string
	Method         string
	Comments       string
	Barcode        string
	CreatedAt      time.Time
	UpdatedAt      time.Time

	// Relationships
	ObservationResult  []ObservationResult  `gorm:"foreignKey:SpecimenID"`
	ObservationRequest []ObservationRequest `gorm:"foreignKey:SpecimenID;->"`
	WorkOrder          WorkOrder            `gorm:"foreignKey:OrderID;->"`
	Patient            Patient              `gorm:"foreignKey:PatientID;->"`
}
