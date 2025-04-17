package entity

import (
	"log"
	"strconv"
	"time"

	"github.com/oibacidem/lims-hl-seven/internal/constant"
)

type ObservationRequest struct {
	ID              int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	TestCode        string    `json:"test_code" gorm:"not null;index:observation_request_uniq,unique,priority:2" validate:"required,observation-type"`
	TestDescription string    `json:"test_description"`
	RequestedDate   time.Time `json:"requested_date"`
	ResultStatus    string    `json:"result_status"`
	SpecimenID      int64     `json:"specimen_id" gorm:"not null;index:observation_request_uniq,unique,priority:1" validate:"required"`
	CreatedAt       time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"not null"`

	TestType TestType `json:"test_type" gorm:"foreignKey:TestCode;references:Code" validate:"-"`
}

func (o ObservationRequest) GetOrderControlNode() string {
	switch o.ResultStatus {
	case string(constant.ResultStatusDelete):
		return string(constant.OrderControlNodeCA)
	default:
		return string(constant.OrderControlNodeNW)
	}
}

type ObservationResultTestsCreate struct {
	TestTypeID int64   `json:"test_type_id" validate:"required"`
	Value      float64 `json:"value" validate:"required"`
}

type ObservationResultCreate struct {
	SpecimenID int64                          `json:"specimen_id" validate:"required"`
	Tests      []ObservationResultTestsCreate `json:"tests" validate:"required"`
}

type ObservationResult struct {
	ID             int64           `json:"id" gorm:"primaryKey;autoIncrement"`
	SpecimenID     int64           `json:"specimen_id"`
	TestCode       string          `json:"code" gorm:"column:code"`
	Description    string          `json:"description"`
	Values         JSONStringArray `json:"values" gorm:"type:json"` // Using JSON for the slice
	Type           string          `json:"type"`
	Unit           string          `json:"unit"`
	ReferenceRange string          `json:"reference_range"`
	Date           time.Time       `json:"date"`
	AbnormalFlag   JSONStringArray `json:"abnormal_flag" gorm:"type:json"` // Using JSON for the slice
	Comments       string          `json:"comments"`
	Picked         bool            `json:"picked" gorm:"not null,default:false"`
	CreatedAt      time.Time       `json:"created_at" gorm:"not null"`
	UpdatedAt      time.Time       `json:"updated_at" gorm:"not null"`

	TestType TestType `json:"test_type" gorm:"foreignKey:TestCode;references:Code" validate:"required"`
}

// GetFirstValue get the first value from the values
// If it needs two or more values, then we need to handle it later
// TODO do we really need more than one values?
func (o ObservationResult) GetFirstValue() float64 {
	if len(o.Values) < 1 {
		log.Printf("values from observation %d is empty or negative", o.ID)
		return 0
	}

	v, err := strconv.ParseFloat(o.Values[0], 64)
	if err != nil {
		log.Printf("parse observation.Values from observation %d failed: %v", o.ID, err)
		return v
	}

	return v
}

type ObservationRequestGetManyRequest struct {
	GetManyRequest
	SpecimenID []int64 `query:"specimen_id"`
}
