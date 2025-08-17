package entity

import (
	"fmt"
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

type WorkOrderVerifiedStatus string

const (
	WorkOrderVerifiedStatusPending  WorkOrderVerifiedStatus = "PENDING"
	WorkOrderVerifiedStatusVerified WorkOrderVerifiedStatus = "VERIFIED"
	WorkOrderVerifiedStatusRejected WorkOrderVerifiedStatus = "REJECTED"
)

type WorkOrderCreateRequest struct {
	PatientID       int64                            `json:"patient_id" validate:"required"`
	TestTypes       []WorkOrderCreateRequestTestType `json:"test_types" validate:"required,min=1"`
	CreatedBy       int64                            `json:"created_by" validate:"required"`
	DoctorIDs       []int64                          `json:"doctor_ids" gorm:"-"`
	AnalyzerIDs     []int64                          `json:"analyzer_ids" gorm:"-"`
	TestTemplateIDs []int64                          `json:"test_template_ids" gorm:"-"`

	Barcode      string `json:"barcode"`
	BarcodeSIMRS string `json:"barcode_simrs"`
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
	//nolint:lll // tag cannot be shorter
	Barcode        string    `json:"barcode" gorm:"column:barcode;type:varchar(255);default:'';index:work_order_barcode,unique"`
	BarcodeSIMRS   string    `json:"barcode_simrs" gorm:"column:barcode_simrs;type:varchar(255);default:''"`
	VerifiedStatus string    `json:"verified_status" gorm:"column:verified_status;type:varchar(255);default:''"`
	CreatedBy      int64     `json:"created_by" gorm:"column:created_by;type:bigint;default:0"`
	LastUpdatedBy  int64     `json:"last_updated_by" gorm:"column:last_updated_by;type:bigint;default:0"`
	CreatedAt      time.Time `json:"created_at" gorm:"index:work_order_created_at"`
	UpdatedAt      time.Time `json:"updated_at" gorm:""`

	DoctorIDs       []int64 `json:"doctor_ids" gorm:"-"`
	AnalyzerIDs     []int64 `json:"analyzer_ids" gorm:"-"`
	TestTemplateIDs []int64 `json:"test_template_ids" gorm:"-"`

	Patient          Patient        `json:"patient" gorm:"foreignKey:PatientID;->" validate:"-"`
	Specimen         []Specimen     `json:"specimen_list,omitempty" gorm:"foreignKey:OrderID;->" validate:"-"`
	Devices          []Device       `json:"devices" gorm:"many2many:work_order_devices;->" validate:"-"`
	CreatedByUser    Admin          `json:"created_by_user" gorm:"foreignKey:CreatedBy;->" validate:"-"`
	LastUpdateByUser Admin          `json:"last_updated_by_user" gorm:"foreignKey:LastUpdatedBy;->" validate:"-"`
	Doctors          []Admin        `json:"doctors" gorm:"many2many:work_order_doctors;->" validate:"-"`
	Analyzers        []Admin        `json:"analyzers" gorm:"many2many:work_order_analyzers;->" validate:"-"`
	TestTemplates    []TestTemplate `json:"test_template" gorm:"many2many:work_order_test_templates;->" validate:"-"`

	TestResult        []TestResult `json:"test_result" gorm:"-"`
	TotalRequest      int64        `json:"total_request" gorm:"-"`
	TotalResultFilled int64        `json:"total_result_filled" gorm:"-"`
	PercentComplete   float64      `json:"percent_complete" gorm:"-"`
	HaveCompleteData  bool         `json:"have_complete_data" gorm:"-"`
}

func (wo *WorkOrder) GetFirstDoctor() Admin {
	if len(wo.Doctors) > 0 {
		return wo.Doctors[0]
	}

	return Admin{}
}

func (wo *WorkOrder) FillData() {
	var doctorIDs []int64
	for _, d := range wo.Doctors {
		doctorIDs = append(doctorIDs, d.ID)
	}

	var analyzerIDs []int64
	for _, a := range wo.Analyzers {
		analyzerIDs = append(analyzerIDs, a.ID)
	}

	var testTemplateIDs []int64
	for _, t := range wo.TestTemplates {
		testTemplateIDs = append(testTemplateIDs, int64(t.ID))
	}

	wo.DoctorIDs = doctorIDs
	wo.AnalyzerIDs = analyzerIDs
	wo.TestTemplateIDs = testTemplateIDs
}

type WorkOrderDoctor struct {
	WorkOrderID int64 `json:"work_order_id" gorm:"primaryKey" validate:"required"`
	AdminID     int64 `json:"admin_id" gorm:"primaryKey" validate:"required"`
}

type WorkOrderAnalyzer struct {
	WorkOrderID int64 `json:"work_order_id" gorm:"primaryKey" validate:"required"`
	AdminID     int64 `json:"admin_id" gorm:"primaryKey" validate:"required"`
}

type WorkOrderTestTemplate struct {
	WorkOrderID    int64 `json:"work_order_id" gorm:"primaryKey" validate:"required"`
	TestTemplateID int64 `json:"test_template_id" gorm:"primaryKey" validate:"required"`
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

	progressWriter chan WorkOrderRunStreamMessage
	patients       []Patient
	device         Device
}

func (w *WorkOrderRunRequest) SetPatients(patients []Patient) {
	w.patients = patients
}

func (w *WorkOrderRunRequest) GetPatients() []Patient {
	return w.patients
}

func (w *WorkOrderRunRequest) SetDevice(device Device) {
	w.device = device
}

func (w *WorkOrderRunRequest) GetDevice() Device {
	return w.device
}

func (w *WorkOrderRunRequest) ProgressWriter() chan WorkOrderRunStreamMessage {
	if w.progressWriter == nil {
		w.progressWriter = make(chan WorkOrderRunStreamMessage)
	}

	return w.progressWriter
}

func (w *WorkOrderRunRequest) SetProgressWriter(progress chan WorkOrderRunStreamMessage) {
	w.progressWriter = progress
}

type WorkOrderGetManyRequest struct {
	GetManyRequest
	BarcodeIds  []int64 `query:"barcode_ids"`
	PatientID   int64   `query:"patient_id"`
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

	ProgressWriter chan WorkOrderRunStreamMessage
}

// WorkOrderRunStreamMessage represents a message sent from the use case to the controller.
type WorkOrderRunStreamMessage struct {
	Percentage float64
	Status     WorkOrderStreamingResponseStatus
	Error      error
	IsDone     bool
}
