package entity

import "time"

type WorkOrderStatus string

const (
	WorkOrderStatusNew       WorkOrderStatus = "NEW"
	WorkOrderStatusPending   WorkOrderStatus = "PENDING"
	WorkOrderCancelled       WorkOrderStatus = "CANCELLED"
	WorkOrderStatusCompleted WorkOrderStatus = "SUCCESS"
)

type WorkOrderCreateRequest struct {
	PatientID int64   `json:"patient_id" validate:"required"`
	TestIDs   []int64 `json:"test_ids" validate:"required"`
}

type WorkOrder struct {
	ID                 int64           `json:"id,omitempty" gorm:"primaryKey;autoIncrement"`
	Status             WorkOrderStatus `json:"status" gorm:"not null" validate:"work-order-status"`
	PatientID          int64           `json:"patient_id" gorm:"type:not null;default:0"`
	DeviceIDDeprecated int64           `json:"device_id" gorm:"column:device_id;type:not null;default:0"`
	CreatedAt          time.Time       `json:"created_at" gorm:"index:work_order_created_at"`
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
	Urgent       bool    `json:"urgent" gorm:"-" validate:"required"`
}

type WorkOrderCancelRequest struct {
	WorkOrderID int64 `json:"work_order_id" gorm:"-" validate:"required"`
	DeviceID    int64 `json:"device_id" gorm:"-" validate:"required"`
}

type WorkOrderGetManyRequest struct {
	GetManyRequest
	PatientIDs  []int64 `query:"patient_ids"`
	SpecimenIDs []int64 `query:"specimen_ids"`
}
