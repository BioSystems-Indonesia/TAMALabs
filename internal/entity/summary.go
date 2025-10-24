package entity

type Summary struct {
	Name  string `json:"name"`
	Total int    `json:"total"`
}

type SummaryResponse struct {
	WorkOrderTrend       []Summary `json:"work_order_trend"`
	AbnormalSummary      []Summary `json:"abnormal_summary"`
	GenderSummary        []Summary `json:"gender_summary"`
	AgeGroup             []Summary `json:"age_group"`
	TopTestOrdered       []Summary `json:"top_test_ordered"`
	TestTypeDistribution []Summary `json:"test_type_distribution"`
}

type SummaryCardResponse struct {
	TotalWorkOrders      int `json:"total_work_orders"`
	CompletedWorkOrders  int `json:"completed_work_orders"`
	PendingWorkOrders    int `json:"pending_work_orders"`
	IncomplateWorkOrders int `json:"incomplate_work_orders"`
	TotalTest            int `json:"total_test"`
	DevicesConnected     int `json:"devices_connected"`
	TotalPatients        int `json:"total_patients"`
	TotalTestParameters  int `json:"total_test_parameters"`
}
