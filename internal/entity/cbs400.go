package entity

import "time"

type CBS400Result struct {
	PatientID  string    `json:"patient_id"`
	TestName   string    `json:"test_name"`
	SampleType string    `json:"sample_type"`
	Value      float64   `json:"value"`
	Unit       string    `json:"unit"`
	Timestamp  time.Time `json:"timestamp"`
}
