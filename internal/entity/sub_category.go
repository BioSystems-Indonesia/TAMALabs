package entity

import "time"

// SubCategory represents a test sub-category
type SubCategory struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"not null;uniqueIndex:idx_sub_categories_name"`
	Code        string    `json:"code" gorm:"not null"`
	Category    string    `json:"category" gorm:"not null"`
	Description string    `json:"description" gorm:"type:text;default:''"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"not null"`
}

// SubCategoryGetManyRequest represents the request to get many sub-categories
type SubCategoryGetManyRequest struct {
	GetManyRequest
	Category string `query:"category"`
}
