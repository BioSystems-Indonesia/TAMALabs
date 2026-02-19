package entity

import (
	"log/slog"
	"strconv"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/constant"
	"gorm.io/gorm"
)

type ObservationRequest struct {
	ID              int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	TestCode        string    `json:"test_code" gorm:"not null" validate:"required,observation-type"`
	TestTypeID      *int      `json:"test_type_id" gorm:"index:observation_request_uniq,unique,priority:1"` // Specific test_type ID selected by user
	TestDescription string    `json:"test_description"`
	RequestedDate   time.Time `json:"requested_date"`
	ResultStatus    string    `json:"result_status"`
	SpecimenID      int64     `json:"specimen_id" gorm:"not null;index:observation_request_uniq,unique,priority:2" validate:"required"`
	PackageID       *int      `json:"package_id" gorm:"default:null"`
	// simrs_index: index value coming from external SIMRS (Nuha). Nullable for backward compatibility.
	SimrsIndex *int      `json:"simrs_index,omitempty" gorm:"column:simrs_index;default:null"`
	CreatedAt  time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"not null"`

	TestType  TestType  `json:"test_type" gorm:"foreignKey:TestTypeID;references:ID;->" validate:"-"`
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
	TestTypeID     *int            `json:"test_type_id" gorm:"index:idx_observation_results_test_type_id"` // Specific test_type ID
	Description    string          `json:"description"`
	Values         JSONStringArray `json:"values" gorm:"type:json"` // Using JSON for the slice
	Type           string          `json:"type"`
	Unit           string          `json:"unit"`
	ReferenceRange string          `json:"reference_range"`
	Date           time.Time       `json:"date"`
	AbnormalFlag   JSONStringArray `json:"abnormal_flag" gorm:"type:json"` // Using JSON for the slice
	Comments       string          `json:"comments"`
	Picked         bool            `json:"picked" gorm:"not null,default:false"`
	CreatedBy      int64           `json:"created_by" gorm:"not null;default:-1"`
	CreatedAt      time.Time       `json:"created_at" gorm:"not null"`
	UpdatedAt      time.Time       `json:"updated_at" gorm:"not null"`

	TestType       TestType `json:"test_type" gorm:"-" validate:"required"` // Changed to gorm:"-" to prevent auto-loading
	CreatedByAdmin Admin    `json:"created_by_admin" gorm:"foreignKey:CreatedBy;references:ID"`

	// Calculated fields (not stored in database)
	EGFR *EGFRCalculation `json:"egfr,omitempty" gorm:"-"`
	// ComputedReferenceRange always uses TestType.GetReferenceRange() instead of stored ReferenceRange
	ComputedReferenceRange string `json:"computed_reference_range" gorm:"-"`
}

func (o *ObservationResult) AfterFind(tx *gorm.DB) error {
	switch o.CreatedBy {
	case int64(constant.CreatedByUnknown):
		o.CreatedByAdmin = Admin{
			ID:       int64(constant.CreatedByUnknown),
			Fullname: "Unknown",
		}
	case int64(constant.CreatedBySystem):
		o.CreatedByAdmin = Admin{
			ID:       int64(constant.CreatedBySystem),
			Fullname: "System",
		}
	}

	// Load TestType - prioritize test_type_id if available (for accurate test selection)
	if o.TestTypeID != nil && *o.TestTypeID > 0 {
		// Use test_type_id for accurate lookup (handles cases where multiple tests share the same code)
		var testType TestType
		err := tx.First(&testType, *o.TestTypeID).Error
		if err == nil {
			o.TestType = testType
		} else {
			slog.Warn("TestType not found by ID, falling back to code lookup",
				"test_type_id", *o.TestTypeID, "test_code", o.TestCode, "observation_result_id", o.ID)
			// Fall back to code lookup if ID not found (backward compatibility)
			o.loadTestTypeByCode(tx)
		}
	} else if o.TestCode != "" {
		// Fallback to code-based lookup for backward compatibility
		o.loadTestTypeByCode(tx)
	}

	// Set computed reference range from TestType
	if o.TestType.ID != 0 {
		o.ComputedReferenceRange = o.TestType.GetReferenceRange()
	} else {
		o.ComputedReferenceRange = o.ReferenceRange // fallback to stored value
	}

	return nil
}

// loadTestTypeByCode loads TestType using test_code with support for alternative_codes
func (o *ObservationResult) loadTestTypeByCode(tx *gorm.DB) {
	var testType TestType
	found := false

	// 1. Try by main code
	err := tx.Where("LOWER(code) = LOWER(?)", o.TestCode).First(&testType).Error
	if err == nil {
		o.TestType = testType
		found = true
	}

	// 2. Try by alias_code
	if !found {
		err = tx.Where("LOWER(alias_code) = LOWER(?) AND alias_code != ''", o.TestCode).First(&testType).Error
		if err == nil {
			o.TestType = testType
			found = true
		}
	}

	// 3. Try by alternative_codes (JSON array search)
	if !found {
		var allTestTypes []TestType
		err = tx.Find(&allTestTypes).Error
		if err == nil {
			for _, tt := range allTestTypes {
				if tt.HasCode(o.TestCode) {
					o.TestType = tt
					found = true
					break
				}
			}
		}
	}

	// 4. Try by name (some devices send human-readable test name instead of code)
	if !found {
		err = tx.Where("LOWER(name) = LOWER(?) AND name != ''", o.TestCode).First(&testType).Error
		if err == nil {
			o.TestType = testType
			found = true
		}
	}

	if !found {
		slog.Warn("TestType not found", "test_code", o.TestCode, "observation_result_id", o.ID)
	}
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
	// Priority 1: Use test_type_id if available
	if o.TestTypeID != nil && *o.TestTypeID > 0 {
		var testType TestType
		err := tx.First(&testType, *o.TestTypeID).Error
		if err == nil {
			o.ReferenceRange = testType.GetReferenceRange()
			slog.Info("generateReferenceRange: set reference range by ID",
				"test_type_id", *o.TestTypeID, "reference_range", o.ReferenceRange)
			return nil
		}
		slog.Warn("generateReferenceRange: TestType not found by ID, falling back to code",
			"test_type_id", *o.TestTypeID, "test_code", o.TestCode)
	}

	// Priority 2: Fall back to test_code lookup
	if o.TestCode == "" {
		return nil
	}

	var testType TestType
	found := false

	// 1. Try by main code (case insensitive)
	err := tx.Where("LOWER(code) = LOWER(?)", o.TestCode).First(&testType).Error
	if err == nil {
		found = true
	}

	// 2. Try by alias_code
	if !found {
		err = tx.Where("LOWER(alias_code) = LOWER(?) AND alias_code != ''", o.TestCode).First(&testType).Error
		if err == nil {
			found = true
		}
	}

	// 3. Try by alternative_codes - load all and check in Go
	if !found {
		var allTestTypes []TestType
		err = tx.Find(&allTestTypes).Error
		if err == nil {
			for _, tt := range allTestTypes {
				if tt.HasCode(o.TestCode) {
					testType = tt
					found = true
					break
				}
			}
		}
	}

	// 4. Try by name (devices sometimes send name instead of code)
	if !found {
		err = tx.Where("LOWER(name) = LOWER(?) AND name != ''", o.TestCode).First(&testType).Error
		if err == nil {
			found = true
		}
	}

	if !found {
		slog.Warn("generateReferenceRange: TestType not found", "test_code", o.TestCode)
		return nil
	}

	o.ReferenceRange = testType.GetReferenceRange()
	slog.Info("generateReferenceRange: set reference range", "test_code", o.TestCode, "reference_range", o.ReferenceRange)
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
		return 0 // Return 0 instead of v (which would be 0 anyway when err != nil)
	}

	return v
}

// GetFirstValueAsString get the first value as string (preserves qualitative values)
func (o ObservationResult) GetFirstValueAsString() string {
	if len(o.Values) < 1 {
		slog.Info("failed to get first values: is empty", "id", o.ID)
		return ""
	}
	return o.Values[0]
}

type ObservationRequestGetManyRequest struct {
	GetManyRequest
	SpecimenID []int64 `query:"specimen_id"`
}
