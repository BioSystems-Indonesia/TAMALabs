package entity

type DeviceType string

const (
	DeviceTypeBA400 DeviceType = "BA400"
	DeviceTypeBA200 DeviceType = "BA200"
	DeviceTypeA15   DeviceType = "A15"
	DeviceTypeOther DeviceType = "OTHER"
)

type Device struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	IPAddress string `json:"ip_address"`
	Port      int    `json:"port"`

	Username string `json:"username"`
	Password string `json:"password"`
	Path     string `json:"path"`
}
