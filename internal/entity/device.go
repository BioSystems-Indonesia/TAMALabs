package entity

import "fmt"

type DeviceType string

const (
	DeviceTypeBA400       DeviceType = "BA400"
	DeviceTypeBA200       DeviceType = "BA200"
	DeviceTypeA15         DeviceType = "A15"
	DeviceTypeOther       DeviceType = "OTHER"
	DeviceTypeAnalyxTria  DeviceType = "ANALYX_TRIA"
	DeviceTypeAnalyxPanca DeviceType = "ANALYX_PANCA"
	DeviceTypeSwelabAlfa  DeviceType = "SWELLAB_ALFA"
	DeviceTypeCoax        DeviceType = "COAX"
)

var TableDeviceType = Tables{
	{ID: string(DeviceTypeBA400), Name: string(DeviceTypeBA400), AdditionalInfo: DeviceCapability{
		CanSend:    true,
		CanReceive: true,
	}},
	{ID: string(DeviceTypeBA200), Name: string(DeviceTypeBA200), AdditionalInfo: DeviceCapability{
		CanSend:    true,
		CanReceive: true,
	}},
	{ID: string(DeviceTypeA15), Name: string(DeviceTypeA15), AdditionalInfo: DeviceCapability{
		CanSend:  true,
		HaveAuth: true,
		HavePath: true,
	}},
	{ID: string(DeviceTypeAnalyxTria), Name: string(DeviceTypeAnalyxTria), AdditionalInfo: DeviceCapability{
		CanReceive: true,
	}},
	{ID: string(DeviceTypeAnalyxPanca), Name: string(DeviceTypeAnalyxPanca), AdditionalInfo: DeviceCapability{
		CanReceive: true,
	}},
	{ID: string(DeviceTypeSwelabAlfa), Name: string(DeviceTypeSwelabAlfa), AdditionalInfo: DeviceCapability{
		CanReceive: true,
	}},
	{ID: string(DeviceTypeOther), Name: string(DeviceTypeOther), AdditionalInfo: DeviceCapability{
		CanSend:    true,
		CanReceive: true,
	}},
	{ID: string(DeviceTypeCoax), Name: string(DeviceTypeCoax), AdditionalInfo: DeviceCapability{
		CanReceive: true,
		UseSerial:  true,
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
