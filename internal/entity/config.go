package entity

type Config struct {
	ID    string `json:"id" gorm:"primaryKey" validate:"required"`
	Value string `json:"value" gorm:"not null" validate:"required"`
}

type ConfigGetManyRequest struct {
	GetManyRequest
}

type ConfigPaginationResponse struct {
	Data []Config `json:"data"`
	PaginationResponse
}
