package entity

import "time"

type Patient struct {
	ID        int64     `json:"id,omitempty" gorm:"primaryKey;autoIncrement"`
	FirstName string    `json:"first_name" gorm:"not null" validate:"required"`
	LastName  string    `json:"last_name" gorm:"not null" validate:"required"`
	Birthdate time.Time `json:"birthdate" gorm:"not null" validate:"required"`
	Sex       string    `json:"sex" gorm:"not null" validate:"required"`
	Location  string    `json:"location" gorm:"not null" validate:"required"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
}

type GetManyRequest struct {
	Order  string `query:"_order"`
	Sort   string `query:"_sort"`
	Start  int    `query:"_start"`
	End    int    `query:"_end"`
	Search string `query:"_search"`
}

func (g GetManyRequest) IsSortDesc() bool {
	if g.Order == "DESC" {
		return true
	}

	return false
}
