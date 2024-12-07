package entity

type ErrorPayload struct {
	Path       string                 `json:"path"`
	StatusCode int                    `json:"status_code"`
	Error      string                 `json:"error"`
	ExtraInfo  map[string]interface{} `json:"extra_info,omitempty"`
}
