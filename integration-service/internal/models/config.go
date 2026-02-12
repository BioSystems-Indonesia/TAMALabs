package models

import "time"

type Config struct {
	LabId        string `json:"lab_id"`
	HospitalName string `json:"hospital_name"`
	Location     string `json:"location"`
}

type LabKey struct {
	Id           int       `json:"id"`
	KeyId        string    `json:"key_id"`
	LabId        string    `json:"lab_id"`
	HospitalName string    `json:"hospital_name"`
	Location     string    `json:"location"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
