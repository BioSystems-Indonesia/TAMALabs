package technomedic

type RequestOrder struct {
	NoOrder      string   `json:"no_order"`
	PatientID    string   `json:"patient_id"`
	ParamRequest []string `json:"param_request"`
	RequestedBy  string   `json:"requested_by"`
	RequestedAt  string   `json:"requested_at"`
}
