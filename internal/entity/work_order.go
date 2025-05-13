package entity

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type WorkOrderStatus string

const (
	WorkOrderStatusNew            WorkOrderStatus = "NEW"
	WorkOrderStatusIncompleteSend WorkOrderStatus = "INCOMPLETE_SEND"
	WorkOrderStatusPending        WorkOrderStatus = "PENDING"
	WorkOrderCancelled            WorkOrderStatus = "CANCELLED"
	WorkOrderStatusCompleted      WorkOrderStatus = "SUCCESS"
)

type WorkOrderCreateRequest struct {
	PatientID int64                            `json:"patient_id" validate:"required"`
	TestTypes []WorkOrderCreateRequestTestType `json:"test_types" validate:"required,min=1"`
	Barcode   string
}

type WorkOrderCreateRequestTestType struct {
	TestTypeID   int64  `json:"test_type_id" validate:"required"`
	TestTypeCode string `json:"test_type_code" validate:"required"`
	SpecimenType string `json:"specimen_type" validate:"required"`
}

type WorkOrder struct {
	ID                 int64           `json:"id,omitempty" gorm:"primaryKey;autoIncrement"`
	Status             WorkOrderStatus `json:"status" gorm:"not null" validate:"work-order-status"`
	PatientID          int64           `json:"patient_id" gorm:"type:not null;default:0"`
	DeviceIDDeprecated int64           `json:"device_id" gorm:"column:device_id;type:not null;default:0"`
	CreatedAt          time.Time       `json:"created_at" gorm:"index:work_order_created_at"`
	Barcode            string          `json:"barcode" gorm:"column:barcode;type:varchar(255);default:'';index:work_order_barcode,unique"`
	UpdatedAt          time.Time       `json:"updated_at" gorm:""`

	Patient  Patient    `json:"patient" gorm:"foreignKey:PatientID;->" validate:"-"`
	Specimen []Specimen `json:"specimen_list,omitempty" gorm:"foreignKey:OrderID;->" validate:"-"`
	// nolint:lll // tag cannot be shorter
	Devices []Device `json:"devices" gorm:"many2many:work_order_devices;->" validate:"-"`

	TestResult        []TestResult `json:"test_result" gorm:"-"`
	TotalRequest      int64        `json:"total_request" gorm:"-"`
	TotalResultFilled int64        `json:"total_result_filled" gorm:"-"`
	PercentComplete   float64      `json:"percent_complete" gorm:"-"`
	HaveCompleteData  bool         `json:"have_complete_data" gorm:"-"`
}

type WorkOrderDevice struct {
	WorkOrderID int64     `json:"work_order_id" gorm:"not null;index:work_order_device_uniq,unique" validate:"required"`
	DeviceID    int64     `json:"device_id" gorm:"not null;index:work_order_device_uniq,unique" validate:"required"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"not null"`
}

type WorkOrderRunRequest struct {
	DeviceID     int64   `json:"device_id" gorm:"-" validate:"required"`
	WorkOrderIDs []int64 `json:"work_order_ids" gorm:"-" validate:"required"`
	Urgent       bool    `json:"urgent" gorm:"-"`
}

type WorkOrderGetManyRequest struct {
	GetManyRequest
	PatientIDs  []int64 `query:"patient_ids"`
	SpecimenIDs []int64 `query:"specimen_ids"`
}

type WorkOrderStreamingResponseStatus string

const (
	WorkOrderStreamingResponseStatusDone       WorkOrderStreamingResponseStatus = "DONE"
	WorkOrderStreamingResponseStatusInProgress WorkOrderStreamingResponseStatus = "IN_PROGRESS"
)

type WorkOrderStreamingResponse string

func NewWorkOrderStreamingResponse(percentage float64, status WorkOrderStreamingResponseStatus) string {
	return fmt.Sprintf("data: percentage=%d&status=%s\n\n", int(percentage), status)
}

type SendPayloadRequest struct {
	Patients []Patient
	Device   Device
	Urgent   bool

	Writer  io.Writer
	Flusher http.Flusher
}
