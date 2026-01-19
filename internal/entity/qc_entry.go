package entity

import "time"

// QCEntry represents QC configuration that must be set before QC process
type QCEntry struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	DeviceID   int       `json:"device_id" gorm:"not null"`
	TestTypeID int       `json:"test_type_id" gorm:"not null"`
	QCLevel    int       `json:"qc_level" gorm:"not null"` // 1, 2, or 3
	LotNumber  string    `json:"lot_number" gorm:"not null"`
	TargetMean float64   `json:"target_mean" gorm:"not null"`
	TargetSD   *float64  `json:"target_sd" gorm:"column:target_sd"`
	Method     string    `json:"method" gorm:"not null;default:'statistic'"` // statistic or manual
	IsActive   bool      `json:"is_active" gorm:"not null;default:true"`
	CreatedBy  string    `json:"created_by" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
	RefMin     float64   `json:"ref_min" gorm:"not null"`
	RefMax     float64   `json:"ref_max" gorm:"not null"`

	// Selected result IDs for each level when multiple results exist for the same day
	Level1SelectedResultID *int `json:"level1_selected_result_id,omitempty" gorm:"column:level1_selected_result_id"`
	Level2SelectedResultID *int `json:"level2_selected_result_id,omitempty" gorm:"column:level2_selected_result_id"`
	Level3SelectedResultID *int `json:"level3_selected_result_id,omitempty" gorm:"column:level3_selected_result_id"`

	// Relations
	Device   *Device    `json:"device,omitempty" gorm:"foreignKey:DeviceID"`
	TestType *TestType  `json:"test_type,omitempty" gorm:"foreignKey:TestTypeID"`
	Results  []QCResult `json:"results,omitempty" gorm:"foreignKey:QCEntryID"`
}

func (QCEntry) TableName() string {
	return "qc_entries"
}

type QCResult struct {
	ID        int `json:"id" gorm:"primaryKey"`
	QCEntryID int `json:"qc_entry_id" gorm:"not null"`

	// Raw analyzer value
	MeasuredValue float64 `json:"measured_value" gorm:"not null"`

	// Reference values (manual = target, statistic = real)
	CalculatedMean float64 `json:"calculated_mean" gorm:"not null"`
	CalculatedSD   float64 `json:"calculated_sd" gorm:"not null"`
	CalculatedCV   float64 `json:"calculated_cv" gorm:"not null"` // SD / Mean * 100

	// Errors (ALWAYS relative to CalculatedMean/SD)
	ErrorSD       float64 `json:"error_sd" gorm:"not null"`       // (Measured - Mean) / SD
	AbsoluteError float64 `json:"absolute_error" gorm:"not null"` // Measured - Mean
	RelativeError float64 `json:"relative_error" gorm:"not null"` // (AbsoluteError / Mean) * 100

	// Levey–Jennings lines
	SD1 float64 `json:"sd_1" gorm:"column:sd_1;not null"` // Mean + 1SD
	SD2 float64 `json:"sd_2" gorm:"column:sd_2;not null"` // Mean + 2SD
	SD3 float64 `json:"sd_3" gorm:"column:sd_3;not null"` // Mean + 3SD

	// QC decision
	Result string `json:"result" gorm:"not null"` // In Control / Warning / Reject

	// Source of reference
	Method string `json:"method" gorm:"not null"` // manual | statistic

	Operator         string    `json:"operator" gorm:"not null"`
	MessageControlID string    `json:"message_control_id"`
	CreatedBy        string    `json:"created_by" gorm:"not null"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`

	// Relations
	QCEntry *QCEntry `json:"qc_entry,omitempty" gorm:"foreignKey:QCEntryID"`

	// Virtual
	ResultCount int `json:"result_count" gorm:"-"`
}

func (QCResult) TableName() string {
	return "qc_results"
}

type GetManyRequestQCEntry struct {
	GetManyRequest
	DeviceID   *int  `query:"device_id"`
	TestTypeID *int  `query:"test_type_id"`
	QCLevel    *int  `query:"qc_level"`
	IsActive   *bool `query:"is_active"`
}

type GetManyRequestQCResult struct {
	GetManyRequest
	QCEntryID  *int    `query:"qc_entry_id"`
	DeviceID   *int    `query:"device_id"`
	TestTypeID *int    `query:"test_type_id"`
	Method     *string `query:"method"`     // statistic or manual
	StartDate  *string `query:"start_date"` // YYYY-MM-DD format
	EndDate    *string `query:"end_date"`   // YYYY-MM-DD format
}

type CreateQCEntryRequest struct {
	DeviceID   int      `json:"device_id" validate:"required"`
	TestTypeID int      `json:"test_type_id" validate:"required"`
	QCLevel    int      `json:"qc_level" validate:"required,min=1,max=3"`
	LotNumber  string   `json:"lot_number" validate:"required"`
	TargetMean float64  `json:"target_mean" validate:"required,gt=0"`
	RefMin     float64  `json:"ref_min"`
	RefMax     float64  `json:"ref_max"`
	TargetSD   *float64 `json:"target_sd"`
	Method     string   `json:"method"`
	CreatedBy  string   `json:"created_by" validate:"required"`
}

type UpdateQCEntryRequest struct {
	LotNumber              string   `json:"lot_number"`
	TargetMean             *float64 `json:"target_mean"`
	TargetSD               *float64 `json:"target_sd"`
	IsActive               *bool    `json:"is_active"`
	Level1SelectedResultID *int     `json:"level1_selected_result_id,omitempty"`
	Level2SelectedResultID *int     `json:"level2_selected_result_id,omitempty"`
	Level3SelectedResultID *int     `json:"level3_selected_result_id,omitempty"`
}

// CreateManualQCResultRequest represents a manual QC input request
type CreateManualQCResultRequest struct {
	DeviceID      int     `json:"device_id" validate:"required"`
	TestTypeID    int     `json:"test_type_id" validate:"required"`
	QCLevel       int     `json:"qc_level" validate:"required,min=1,max=3"`
	MeasuredValue float64 `json:"measured_value" validate:"required"`
}

// QCSummary represents aggregated QC statistics for a device
type QCSummary struct {
	DeviceID       int        `json:"device_id"`
	QCTodayStatus  string     `json:"qc_today_status"`        // "Done", "Not Done", or "Partial"
	TotalQC        int        `json:"total_qc"`               // Total QC results count
	QCThisMonth    int        `json:"qc_this_month"`          // QC count for current month
	LastQCDate     *time.Time `json:"last_qc_date,omitempty"` // Date of last QC result
	LastQCStatus   string     `json:"last_qc_status"`         // "Normal", "Warning", or "Error"
	Level1Complete bool       `json:"level_1_complete"`       // Has active entry for level 1
	Level2Complete bool       `json:"level_2_complete"`       // Has active entry for level 2
	Level3Complete bool       `json:"level_3_complete"`       // Has active entry for level 3
	// Per-level QC done today
	Level1Today bool `json:"level_1_today"`
	Level2Today bool `json:"level_2_today"`
	Level3Today bool `json:"level_3_today"`
}
