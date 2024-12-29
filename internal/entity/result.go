package entity

import "time"

type Result struct {
	Date        time.Time    `json:"date"`
	Barcode     string       `json:"barcode"`
	PatientName string       `json:"patient_name"`
	PatientID   int64        `json:"patient_id"`
	Detail      ResultDetail `json:"detail,omitempty"`
}

type ResultDetail struct {
	Hematology   []ResultTest `json:"hematology"`
	Biochemistry []ResultTest `json:"biochemistry"`
	Observation  []ResultTest `json:"observation"`
}

type ResultTest struct {
	Test           string         `json:"test"`
	Result         string         `json:"result"`
	Unit           string         `json:"unit"`
	Category       string         `json:"category"`
	Abnormal       AbnormalResult `json:"abnormal"`
	ReferenceRange string         `json:"reference_range"`
}

type AbnormalResult int32

const (
	NormalResult AbnormalResult = 0
	HighResult   AbnormalResult = 1
	LowResult    AbnormalResult = 2
)

type ResultGetManyRequest struct {
	GetManyRequest
}
