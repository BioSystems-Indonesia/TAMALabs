package entity

import "time"

type BTSResponse struct {
	Message string `json:"message"`
}

type BTSMeasurementData struct {
	MeasureType string  `json:"measureType"`
	ID          string  `json:"id"`
	Analyte     string  `json:"analyte"`
	Result      float64 `json:"result"`
	Units       string  `json:"units"`
	Absorbance  float64 `json:"absorbance"`
	Date        string  `json:"date"`
	Time        string  `json:"time"`
	Status      string  `json:"status,omitempty"`
	LowerRange  string  `json:"lowerRange,omitempty"`
	UpperRange  string  `json:"upperRange,omitempty"`
}
type BTSRequest struct {
	Command BTSCommand `json:"command"`
}

type BTSCommand struct {
	Type    string     `json:"type"`
	Payload BTSPayload `json:"payload"`
}

type BTSPayload struct {
	OperationID   string           `json:"operationId"`
	Filters       BTSFilters       `json:"filters"`
	ExportOptions BTSExportOptions `json:"exportOptions"`
}

type BTSExportOptions struct {
	FileName        string `json:"fileName"`
	Format          string `json:"format"`
	Target          string `json:"target"`
	DeleteAfterDump bool   `json:"deleteAfterDump"`
}

type BTSFilters struct {
	ReadingType []string    `json:"readingType"`
	SampleID    interface{} `json:"sampleId"`
	TechniqueID interface{} `json:"techniqueId"`
	StartDate   time.Time   `json:"startDate"`
	EndDate     time.Time   `json:"endDate"`
	TimePeriod  string      `json:"timePeriod"`
}
