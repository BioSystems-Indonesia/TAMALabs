package entity

import "time"

type AbbottResult struct {
	PatientID string    `json:"patient_id"`
	SampleID  string    `json:"sample_id"`
	TestName  string    `json:"test_name"`
	Value     string    `json:"value"`
	Unit      string    `json:"unit"`
	Flag      string    `json:"flag"`
	RefMin    string    `json:"ref_min"`
	RefMax    string    `json:"ref_max"`
	Timestamp time.Time `json:"timestamp"`
}

type AbbottMessage struct {
	DeviceInfo  AbbottDeviceInfo   `json:"device_info"`
	SampleInfo  AbbottSampleInfo   `json:"sample_info"`
	TestResults []AbbottTestResult `json:"test_results"`
	Timestamp   time.Time          `json:"timestamp"`
}

type AbbottDeviceInfo struct {
	DeviceName string `json:"device_name"`
	Unit       string `json:"unit"`
	Operator   string `json:"operator"`
	Mode       string `json:"mode"`
}

type AbbottSampleInfo struct {
	SampleID    string `json:"sample_id"`
	PatientID   string `json:"patient_id"`
	PatientName string `json:"patient_name"`
	Type        string `json:"type"`
	Date        string `json:"date"`
	Time        string `json:"time"`
	SeqNumber   string `json:"seq_number"`
}

type AbbottTestResult struct {
	TestCode string `json:"test_code"`
	Value    string `json:"value"`
	Unit     string `json:"unit"`
	Flag     string `json:"flag"`
	RefMin   string `json:"ref_min"`
	RefMax   string `json:"ref_max"`
}
