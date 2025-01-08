package entity

type Config struct {
	Key   string `json:"key" gorm:"primaryKey" validate:"required"`
	Value string `json:"value" gorm:"not null" validate:"required"`
}
