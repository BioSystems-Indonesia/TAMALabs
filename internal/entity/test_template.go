package entity

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type TestTemplate struct {
	ID                int       `json:"id" gorm:"primaryKey"`
	Name              string    `json:"name" gorm:"not null" validate:"required"`
	Description       string    `json:"description" gorm:"not null"`
	TestTypesString   string    `json:"-" gorm:"not null;column:test_types;default:{}"`
	DoctorIDsString   string    `json:"-" gorm:"not null;column:doctor_ids;default:[]"`
	AnalyzerIDsString string    `json:"-" gorm:"not null;column:analyzer_ids;default:[]"`
	CreatedBy         int64     `json:"created_by" gorm:"not null;default:0"`
	LastUpdatedBy     int64     `json:"last_updated_by" gorm:"not null;default:0"`
	CreatedAt         time.Time `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`

	RequestTestTypes  []WorkOrderCreateRequestTestType `json:"test_types" gorm:"-" validate:"required"`
	DoctorIDs         []int64                          `json:"doctor_ids" gorm:"-"`
	AnalyzerIDs       []int64                          `json:"analyzer_ids" gorm:"-"`
	CreatedByUser     Admin                            `json:"created_by_user" gorm:"foreignKey:CreatedBy;references:ID;->" validate:"-"`
	LastUpdatedByUser Admin                            `json:"last_updated_by_user" gorm:"foreignKey:LastUpdatedBy;references:ID;->" validate:"-"`
}

func (t *TestTemplate) BeforeCreate(tx *gorm.DB) error {
	testTypeByte, err := json.Marshal(t.RequestTestTypes)
	if err != nil {
		return fmt.Errorf("error marshalling test_types: %w", err)
	}
	t.TestTypesString = string(testTypeByte)

	doctorIDsByte, err := json.Marshal(t.DoctorIDs)
	if err != nil {
		return fmt.Errorf("error marshalling doctor_ids: %w", err)
	}
	t.DoctorIDsString = string(doctorIDsByte)

	analyzerIDsByte, err := json.Marshal(t.AnalyzerIDs)
	if err != nil {
		return fmt.Errorf("error marshalling analyzer_ids: %w", err)
	}
	t.AnalyzerIDsString = string(analyzerIDsByte)

	return nil
}

func (t *TestTemplate) BeforeUpdate(tx *gorm.DB) error {
	testTypeByte, err := json.Marshal(t.RequestTestTypes)
	if err != nil {
		return fmt.Errorf("error marshalling test_types: %w", err)
	}
	t.TestTypesString = string(testTypeByte)

	doctorIDsByte, err := json.Marshal(t.DoctorIDs)
	if err != nil {
		return fmt.Errorf("error marshalling doctor_ids: %w", err)
	}
	t.DoctorIDsString = string(doctorIDsByte)

	analyzerIDsByte, err := json.Marshal(t.AnalyzerIDs)
	if err != nil {
		return fmt.Errorf("error marshalling analyzer_ids: %w", err)
	}
	t.AnalyzerIDsString = string(analyzerIDsByte)

	return nil
}

func (t *TestTemplate) AfterFind(tx *gorm.DB) error {
	var testTypes []WorkOrderCreateRequestTestType
	err := json.Unmarshal([]byte(t.TestTypesString), &testTypes)
	if err != nil {
		return fmt.Errorf("error unmarshalling test_types: %w", err)
	}
	t.RequestTestTypes = testTypes

	var doctorIDs []int64
	err = json.Unmarshal([]byte(t.DoctorIDsString), &doctorIDs)
	if err != nil {
		return fmt.Errorf("error unmarshalling doctor_ids: %w", err)
	}
	t.DoctorIDs = doctorIDs

	var analyzerIDs []int64
	err = json.Unmarshal([]byte(t.AnalyzerIDsString), &analyzerIDs)
	if err != nil {
		return fmt.Errorf("error unmarshalling analyzer_ids: %w", err)
	}
	t.AnalyzerIDs = analyzerIDs

	return nil
}

type TestTemplateTestType struct {
	TestTemplateID int       `json:"test_template_id" gorm:"primaryKey"`
	TestTypeID     int       `json:"test_type_id" gorm:"primaryKey"`
	CreatedAt      time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"not null"`
}

type TestTemplateGetManyRequest struct {
	GetManyRequest
}
