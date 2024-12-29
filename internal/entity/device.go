package entity

type Device struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	IPAddress string `json:"ip_address"`
	Port      int    `json:"port"`
}
