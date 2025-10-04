package simrsuc

// Request represents a SIMRS lab request structure
type Request struct {
	NoOrder      string   `json:"no_order"`
	PatientID    string   `json:"patient_id"`
	ParamRequest []string `json:"param_request"`
	RequestedBy  string   `json:"requested_by"`
	RequestedAt  string   `json:"requested_at"`
}

// Response represents a SIMRS lab result response structure
type Response struct {
	NoOrder string               `json:"no_order"`
	Results []ResponseResultTest `json:"results"`
}

type ResponseResultTest struct {
	ParamCode   string `json:"param_code"`
	ResultValue string `json:"result_value"`
	Unit        string `json:"unit"`
	RefRange    string `json:"ref_range"`
	Flag        string `json:"flag"`
}
