package entity

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type TestTemplate struct {
	ID               int                              `json:"id" gorm:"primaryKey"`
	Name             string                           `json:"name" gorm:"not null" validate:"required"`
	Description      string                           `json:"description" gorm:"not null"`
	TestTypesString  string                           `json:"-" gorm:"not null;column:test_types"`
	RequestTestTypes []WorkOrderCreateRequestTestType `json:"test_types" gorm:"-" validate:"required"`
}

func (t *TestTemplate) BeforeCreate(tx *gorm.DB) error {
	testTypeByte, err := json.Marshal(t.RequestTestTypes)
	if err != nil {
		return err
	}

	t.TestTypesString = string(testTypeByte)

	return nil
}

func (t *TestTemplate) BeforeUpdate(tx *gorm.DB) error {
	testTypeByte, err := json.Marshal(t.RequestTestTypes)
	if err != nil {
		return err
	}

	t.TestTypesString = string(testTypeByte)

	return nil
}

func (t *TestTemplate) AfterFind(tx *gorm.DB) error {
	var testTypes []WorkOrderCreateRequestTestType
	err := json.Unmarshal([]byte(t.TestTypesString), &testTypes)
	if err != nil {
		return err
	}

	t.RequestTestTypes = testTypes

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
