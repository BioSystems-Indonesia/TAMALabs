package entity

import "fmt"

type DeviceType string

const (
	DeviceTypeA15         DeviceType = "A15"
	DeviceTypeSwelabAlfa  DeviceType = "SWELAB_ALFA_PLUS_SATANDARD"
	DeviceTypeSwelabBasic DeviceType = "SWELAB_ALFA_PLUS_BASIC"
)

func (d DeviceType) String() string {
	return string(d)
}

var TableDeviceType = Tables{
	{ID: string(DeviceTypeA15), Name: string(DeviceTypeA15), AdditionalInfo: DeviceCapability{
		HavePath:   true,
		CanReceive: true,
		CanSend:    true,
		HaveAuth:   true,
	}},
	{ID: string(DeviceTypeSwelabAlfa), Name: string(DeviceTypeSwelabAlfa), AdditionalInfo: DeviceCapability{
		CanReceive: true,
	}},
	{ID: string(DeviceTypeSwelabBasic), Name: string(DeviceTypeSwelabBasic), AdditionalInfo: DeviceCapability{
		CanReceive: true,
	}},
}

type DeviceCapability struct {
	CanSend    bool `json:"can_send"`
	CanReceive bool `json:"can_receive"`
	HaveAuth   bool `json:"have_authentication"`
	HavePath   bool `json:"have_path"`
	UseSerial  bool `json:"use_serial"`
}

type Device struct {
	ID          int        `json:"id" gorm:"primaryKey"`
	Name        string     `json:"name"`
	Type        DeviceType `json:"type"`
	IPAddress   string     `json:"ip_address"`
	SendPort    string     `json:"send_port"`
	ReceivePort string     `json:"receive_port" gorm:"unique"`

	Username string `json:"username"`
	Password string `json:"password"`
	Path     string `json:"path"`

	BaudRate int `json:"baud_rate" gorm:"default:0"`
}

type GetManyRequestDevice struct {
	GetManyRequest
}

type DeviceConnectionStatus string

const (
	DeviceConnectionStatusConnected    DeviceConnectionStatus = "connected"
	DeviceConnectionStatusStandby      DeviceConnectionStatus = "standby"
	DeviceConnectionStatusNotSupported DeviceConnectionStatus = "not_supported"
	DeviceConnectionStatusDisconnected DeviceConnectionStatus = "disconnected"
)

type DeviceConnectionMessage string

func NewDeviceConnectionMessage(resp DeviceConnectionResponse) DeviceConnectionMessage {
	return DeviceConnectionMessage(fmt.Sprintf(
		"data: device_id=%d&sender_message=%s&sender_status=%s&receiver_status=%s&receiver_message=%s\n\n",
		resp.DeviceID,
		resp.Sender.Message,
		resp.Sender.Status,
		resp.Receiver.Status,
		resp.Receiver.Message,
	))
}

type DeviceConnectionResponse struct {
	DeviceID int                            `json:"device_id"`
	Sender   DeviceConnectionStatusResponse `json:"sender"`
	Receiver DeviceConnectionStatusResponse `json:"receiver"`
}

type DeviceConnectionStatusResponse struct {
	Message string                 `json:"message"`
	Status  DeviceConnectionStatus `json:"status"`
}

type DeviceConnectionRequest struct {
	DeviceIDs      []int `query:"device_ids" validate:"required,min=1,max=100"`
	TimeoutSeconds int   `query:"timeout_seconds"`
}
