package dto

import "time"

type ObservationResultRequest struct {
	Id             string                  `json:"id"`
	LabID          string                  `json:"lab_id"`
	OrderId        string                  `json:"order_id"`
	CollectionDate string                  `json:"collection_date"`
	ComplatedAt    time.Time               `json:"complated_at"`
	IsVerified     bool                    `json:"is_verified"`
	Patient        PatientRequest          `json:"patient"`
	Items          []ObservationResultItem `json:"items"`
}

type ObservationResultItem struct {
	Id       string `json:"id"`
	Category string `json:"category"`
	Code     string `json:"code"`
	Value    string `json:"value"`
	Flag     string `json:"flag"`
	AddedBy  string `json:"added_by"`
}

type VerifyObservationResultRequest struct {
	Event     string    `json:"event"`
	LabId     string    `json:"lab_id"`
	ResultId  string    `json:"result_id"`
	Status    string    `json:"status"`
	DoctorId  string    `json:"doctor_id"`
	Timestamp time.Time `json:"timestamp"`
}
