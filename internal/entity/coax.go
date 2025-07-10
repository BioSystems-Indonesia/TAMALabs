package entity

type CoaxTestResult struct {
	RecordType string   `json:"record_type"`
	DeviceID   string   `json:"device_id"`
	Status     string   `json:"status"`
	Date       string   `json:"date"`
	Time       string   `json:"time"`
	TestType   string   `json:"test_type"`
	TestName   string   `json:"test_name"`
	Value      string   `json:"value"`
	Unit       string   `json:"unit"`
	Reference  string   `json:"reference"`
	Flags      string   `json:"flags"`
	Extra      []string `json:"extra"` // To handle additional fields
}
