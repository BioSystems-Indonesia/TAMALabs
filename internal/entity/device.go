package entity

type Device struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	Name      string `json:"name"`
	IPAddress string `json:"ip_address"`
}
