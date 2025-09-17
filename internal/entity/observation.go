package entity

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/oibacidem/lims-hl-seven/internal/constant"
	"gorm.io/gorm"
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

	TestType  TestType  `json:"test_type" gorm:"foreignKey:TestCode;references:Code;->" validate:"-"`
	WorkOrder WorkOrder `json:"work_order" gorm:"-" validate:"-"`
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

type EGFRCalculation struct {
	Value    float64 `json:"value"`
	Formula  string  `json:"formula"`
	Unit     string  `json:"unit"`
	Category string  `json:"category"`
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

	// Calculated fields (not stored in database)
	EGFR *EGFRCalculation `json:"egfr,omitempty" gorm:"-"`
}

// BeforeCreate is called before creating ObservationResult
func (o *ObservationResult) BeforeCreate(tx *gorm.DB) error {
	return o.generateReferenceRange(tx)
}

// BeforeUpdate is called before updating ObservationResult
func (o *ObservationResult) BeforeUpdate(tx *gorm.DB) error {
	return o.generateReferenceRange(tx)
}

// generateReferenceRange generates reference range based on TestType
func (o *ObservationResult) generateReferenceRange(tx *gorm.DB) error {
	if o.TestCode == "" {
		return nil
	}

	var testType TestType
	err := tx.Where("code = ?", o.TestCode).First(&testType).Error
	if err != nil {
		// Try to find by alias_code if not found by code
		err = tx.Where("alias_code = ? AND alias_code != ''", o.TestCode).First(&testType).Error
		if err != nil {
			return nil // Skip if test type not found
		}
	}

	decimal := testType.Decimal
	if decimal < 0 {
		decimal = 0
	}

	o.ReferenceRange = fmt.Sprintf("%.*f - %.*f", decimal, testType.LowRefRange, decimal, testType.HighRefRange)
	return nil
}

// GetFirstValue get the first value from the values
// If it needs two or more values, then we need to handle it later
// TODO do we really need more than one values?
func (o ObservationResult) GetFirstValue() float64 {
	if len(o.Values) < 1 {
		slog.Info("failed to get first values: is empty", "id", o.ID)
		return 0
	}

	v, err := strconv.ParseFloat(o.Values[0], 64)
	if err != nil {
		slog.Warn("failed to parse observation.Values from observation", "id", o.ID, "error", err)
		return v
	}

	return v
}

type ObservationRequestGetManyRequest struct {
	GetManyRequest
	SpecimenID []int64 `query:"specimen_id"`
}
