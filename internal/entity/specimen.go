package entity

import (
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/oibacidem/lims-hl-seven/internal/util"
)

type Specimen struct {
	ID             int       `json:"id" gorm:"primaryKey;autoIncrement"`
	HL7ID          string    `json:"specimen_hl7_id" gorm:"not null"`                                                              // SPM-2
	PatientID      int       `json:"patient_id" gorm:"not null;index:specimen_uniq,unique,priority:2" validate:"required"`         // Foreign key linking to Patient
	OrderID        int       `json:"order_id" gorm:"not null;index:specimen_uniq,unique,priority:1" validate:"required"`           // Foreign key linking to WorkOrder
	Type           string    `json:"type" gorm:"not null;index:specimen_uniq,unique,priority:3" validate:"required,specimen-type"` // SPM-4
	CollectionDate string    `json:"collection_date" gorm:"not null"`                                                              // SPM-17
	ReceivedDate   time.Time `json:"received_date" gorm:"not null"`
	Source         string    `json:"source" gorm:"not null"`    // SPM-8
	Condition      string    `json:"condition" gorm:"not null"` // SPM-25
	Method         string    `json:"method" gorm:"not null"`    // SPM-10
	Comments       string    `json:"comments" gorm:"not null"`  // SPM-26
	Barcode        string    `json:"barcode" gorm:"not null"`
	CreatedAt      time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"not null"`

	// Relationships
	Observation        []Observation        `json:"observation" gorm:"-" validate:"-"`
	ObservationRequest []ObservationRequest `json:"observation_requests" gorm:"foreignKey:SpecimenID;->" validate:"-"`
	WorkOrder          WorkOrder            `json:"work_order" gorm:"foreignKey:OrderID;->" validate:"-"`
	Patient            Patient              `json:"patient" gorm:"foreignKey:PatientID;->" validate:"-"`
}

func GenerateBarcode() string {
	randomDigits, err := util.GenerateRandomDigits(4)
	if err != nil {
		log.Errorj(map[string]interface{}{
			"message": "error generating barcode",
			"error":   err,
		})
	}

	return fmt.Sprintf("%s%s", time.Now().Format("20060102"), randomDigits)
}

type SpecimenGetManyRequest struct {
	GetManyRequest
	PatientID int64 `query:"patient_id"`
}
