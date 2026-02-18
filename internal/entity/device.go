package entity

import "fmt"

type DeviceType string

const (
	// Multisera Group
	DeviceTypeBA400            DeviceType = "BA400"
	DeviceTypeBA200            DeviceType = "BA200"
	DeviceTypeA15              DeviceType = "A15"
	DeviceTypeOther            DeviceType = "OTHER"
	DeviceTypeAnalyxTria       DeviceType = "ANALYX_TRIA"
	DeviceTypeAnalyxPanca      DeviceType = "ANALYX_PANCA"
	DeviceTypeSwelabAlfa       DeviceType = "SWELAB_ALFA"
	DeviceTypeSwelabBasic      DeviceType = "SWELAB_BASIC"
	DeviceTypeSwelabLumi       DeviceType = "SWELAB_LUMI"
	DeviceTypeCoax             DeviceType = "COAX"
	DeviceTypeNeomedicaNCC3300 DeviceType = "NEOMEDICA_NCC_3300"
	DeviceTypeNeomedicaNCC61   DeviceType = "NEOMEDICA_NCC_61"
	DeviceTypeAlifax           DeviceType = "ALIFAX"
	DeviceTypeBTS              DeviceType = "BTS"
	DeviceTypeDiestro          DeviceType = "DIESTRO"
	DeviceTypeAbbott           DeviceType = "ABBOTT"
	DeviceTypeResponse911      DeviceType = "RESPONSE_911"
	DeviceTypeEdanI15          DeviceType = "EDAN_I15"
	DeviceTypeEdanH30          DeviceType = "EDAN_H30"

	// Others
	DeviceTypeWondfo     DeviceType = "Wondfo"
	DeviceTypeCBS400     DeviceType = "CBS400"
	DeviceTypeVerifyU120 DeviceType = "VerifyU120"
)

func (d DeviceType) String() string {
	return string(d)
}

var TableDeviceType = Tables{
	{ID: string(DeviceTypeBA200), Name: string(DeviceTypeBA200), AdditionalInfo: DeviceCapability{
		CanSend:    true,
		CanReceive: true,
	}},
	{ID: string(DeviceTypeBA400), Name: string(DeviceTypeBA400), AdditionalInfo: DeviceCapability{
		CanSend:    true,
		CanReceive: true,
	}},
	{ID: string(DeviceTypeA15), Name: string(DeviceTypeA15), AdditionalInfo: DeviceCapability{
		CanReceive: true,
		CanSend:    true,
	}},
	{ID: string(DeviceTypeAnalyxTria), Name: string(DeviceTypeAnalyxTria), AdditionalInfo: DeviceCapability{
		CanReceive: true,
	}},
	{ID: string(DeviceTypeAnalyxPanca), Name: string(DeviceTypeAnalyxPanca), AdditionalInfo: DeviceCapability{
		CanReceive: true,
	}},

	{ID: string(DeviceTypeOther), Name: string(DeviceTypeOther), AdditionalInfo: DeviceCapability{
		CanSend:  true,
		HaveAuth: true,
	}},
	{ID: string(DeviceTypeSwelabAlfa), Name: string(DeviceTypeSwelabAlfa), AdditionalInfo: DeviceCapability{
		CanReceive: true,
	}},
	{ID: string(DeviceTypeSwelabBasic), Name: string(DeviceTypeSwelabBasic), AdditionalInfo: DeviceCapability{
		CanReceive: true,
	}},
	{ID: string(DeviceTypeSwelabLumi), Name: string(DeviceTypeSwelabLumi), AdditionalInfo: DeviceCapability{
		CanReceive: true,
	}},
	{ID: string(DeviceTypeAlifax), Name: string(DeviceTypeAlifax), AdditionalInfo: DeviceCapability{
		CanReceive: true,
		UseSerial:  true,
	}},
	{ID: string(DeviceTypeNeomedicaNCC61), Name: string(DeviceTypeNeomedicaNCC61), AdditionalInfo: DeviceCapability{
		CanReceive: true,
	}},
	{ID: string(DeviceTypeNeomedicaNCC3300), Name: string(DeviceTypeNeomedicaNCC3300), AdditionalInfo: DeviceCapability{
		CanReceive: true,
		UseSerial:  true,
	}},
	{ID: string(DeviceTypeBTS), Name: string(DeviceTypeBTS), AdditionalInfo: DeviceCapability{
		CanSend: true,
	}},
	{ID: string(DeviceTypeCoax), Name: string(DeviceTypeCoax), AdditionalInfo: DeviceCapability{
		CanReceive: true,
		UseSerial:  true,
	}},
	{ID: string(DeviceTypeDiestro), Name: string(DeviceTypeDiestro), AdditionalInfo: DeviceCapability{
		CanReceive: true,
		UseSerial:  true,
	}},
	{ID: string(DeviceTypeWondfo), Name: string(DeviceTypeWondfo), AdditionalInfo: DeviceCapability{
		CanReceive: true,
	}},
	{ID: string(DeviceTypeCBS400), Name: string(DeviceTypeCBS400), AdditionalInfo: DeviceCapability{
		CanReceive: true,
		UseSerial:  true,
	}},
	{ID: string(DeviceTypeVerifyU120), Name: string(DeviceTypeVerifyU120), AdditionalInfo: DeviceCapability{
		CanReceive: true,
		UseSerial:  true,
	}},
	{ID: string(DeviceTypeAbbott), Name: string(DeviceTypeAbbott), AdditionalInfo: DeviceCapability{
		CanSend:    true,
		CanReceive: true,
	}},
	{ID: string(DeviceTypeResponse911), Name: string(DeviceTypeResponse911), AdditionalInfo: DeviceCapability{
		CanReceive: true,
	}},
	{ID: string(DeviceTypeEdanI15), Name: string(DeviceTypeEdanI15), AdditionalInfo: DeviceCapability{
		CanReceive: true,
	}},
	{ID: string(DeviceTypeEdanH30), Name: string(DeviceTypeEdanH30), AdditionalInfo: DeviceCapability{
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
	Type []DeviceType `query:"types"`
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
