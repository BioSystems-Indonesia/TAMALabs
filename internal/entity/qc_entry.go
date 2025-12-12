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

	// Relations
	Device   *Device    `json:"device,omitempty" gorm:"foreignKey:DeviceID"`
	TestType *TestType  `json:"test_type,omitempty" gorm:"foreignKey:TestTypeID"`
	Results  []QCResult `json:"results,omitempty" gorm:"foreignKey:QCEntryID"`
}

func (QCEntry) TableName() string {
	return "qc_entries"
}

// QCResult represents actual measurement result from analyzer
type QCResult struct {
	ID               int       `json:"id" gorm:"primaryKey"`
	QCEntryID        int       `json:"qc_entry_id" gorm:"not null"`
	MeasuredValue    float64   `json:"measured_value" gorm:"not null"`
	CalculatedMean   float64   `json:"calculated_mean" gorm:"not null"`
	CalculatedSD     float64   `json:"calculated_sd" gorm:"not null"`
	CalculatedCV     float64   `json:"calculated_cv" gorm:"not null"`                        // SD/Mean
	ErrorSD          float64   `json:"error_sd" gorm:"column:error_sd;not null"`             // (MeasuredValue - TargetMean) / TargetSD for Levey-Jennings chart
	AbsoluteError    float64   `json:"absolute_error" gorm:"column:absolute_error;not null"` // MeasuredValue - TargetMean
	RelativeError    float64   `json:"relative_error" gorm:"column:relative_error;not null"` // (AbsoluteError / TargetMean) * 100
	SD1              float64   `json:"sd_1" gorm:"column:sd_1;not null"`                     // Mean + (1*SD)
	SD2              float64   `json:"sd_2" gorm:"column:sd_2;not null"`                     // Mean + (2*SD)
	SD3              float64   `json:"sd_3" gorm:"column:sd_3;not null"`                     // Mean + (3*SD)
	Result           string    `json:"result" gorm:"not null"`                               // Pass or Fail
	Method           string    `json:"method" gorm:"not null;default:'statistic'"`           // statistic or manual
	Operator         string    `json:"operator" gorm:"not null"`
	MessageControlID string    `json:"message_control_id" gorm:"column:message_control_id"`
	CreatedAt        time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Relations
	QCEntry     *QCEntry `json:"qc_entry,omitempty" gorm:"foreignKey:QCEntryID"`
	ResultCount int      `json:"result_count" gorm:"-"` // Count of results for this entry (not stored in DB)
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
	QCEntryID *int    `query:"qc_entry_id"`
	DeviceID  *int    `query:"device_id"`
	Method    *string `query:"method"` // statistic or manual
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
	Method     string   `json:"method" validate:"required,oneof=statistic manual"`
	CreatedBy  string   `json:"created_by" validate:"required"`
}

type UpdateQCEntryRequest struct {
	LotNumber  string   `json:"lot_number"`
	TargetMean *float64 `json:"target_mean"`
	TargetSD   *float64 `json:"target_sd"`
	IsActive   *bool    `json:"is_active"`
}

// QCSummary represents aggregated QC statistics for a device
type QCSummary struct {
	DeviceID       int       `json:"device_id"`
	QCTodayStatus  string    `json:"qc_today_status"`  // "Done", "Not Done", or "Partial"
	TotalQC        int       `json:"total_qc"`         // Total QC results count
	QCThisMonth    int       `json:"qc_this_month"`    // QC count for current month
	LastQCDate     time.Time `json:"last_qc_date"`     // Date of last QC result
	LastQCStatus   string    `json:"last_qc_status"`   // "Normal", "Warning", or "Error"
	Level1Complete bool      `json:"level_1_complete"` // Has active entry for level 1
	Level2Complete bool      `json:"level_2_complete"` // Has active entry for level 2
	Level3Complete bool      `json:"level_3_complete"` // Has active entry for level 3
}
