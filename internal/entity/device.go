package entity

import "fmt"

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

type GetManyRequestDevice struct {
	GetManyRequest
}

type DeviceConnectionStatus string

const (
	DeviceConnectionStatusConnected    DeviceConnectionStatus = "connected"
	DeviceConnectionStatusNotSupported DeviceConnectionStatus = "not_supported"
	DeviceConnectionStatusDisconnected DeviceConnectionStatus = "disconnected"
)

type DeviceConnectionMessage string

func NewDeviceConnectionMessage(deviceID int, message string, status DeviceConnectionStatus) DeviceConnectionMessage {
	return DeviceConnectionMessage(fmt.Sprintf("data: device_id=%d&message=%s&status=%s\n\n", deviceID, message, status))
}

type DeviceConnectionResponse struct {
	DeviceID int                    `json:"device_id"`
	Message  string                 `json:"message"`
	Status   DeviceConnectionStatus `json:"status"`
}

type DeviceConnectionRequest struct {
	DeviceIDs      []int `query:"device_ids" validate:"required,min=1,max=100"`
	TimeoutSeconds int   `query:"timeout_seconds"`
}
