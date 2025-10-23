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
