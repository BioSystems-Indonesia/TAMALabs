package entity

import "time"

type VerifyResult struct {
	PatientID  string    `json:"patient_id"`
	TestName   string    `json:"test_name"`
	SampleType string    `json:"sample_type"`
	Value      float64   `json:"value"`
	ValueStr   string    `json:"value_str"`
	Unit       string    `json:"unit"`
	Timestamp  time.Time `json:"timestamp"`
}
