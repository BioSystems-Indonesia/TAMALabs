package entity

import "time"

type TestTemplate struct {
	ID          int    `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"not null" validate:"required"`
	Description string `json:"description" gorm:"not null"`
	TestTypeID  []int  `json:"test_type_id" gorm:"-"`

	TestType []TestType `json:"test_type,omitempty" gorm:"many2many:test_template_test_types;->"`
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

type TestTemplatePaginationResponse struct {
	TestTemplates []TestTemplate `json:"test_templates"`
	PaginationResponse
}
