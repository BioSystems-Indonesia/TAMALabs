package entity

// TechnoMedicOrderRequest represents the request to create an order from TechnoMedic
type TechnoMedicOrderRequest struct {
	NoOrder            string                    `json:"no_order" validate:"required"`
	Patient            TechnoMedicPatientRequest `json:"patient" validate:"required"`
	ParamRequest       []string                  `json:"param_request"`        // Test type codes (e.g., ["HB", "WBC"])
	SubCategoryRequest []string                  `json:"sub_category_request"` // Sub-category names (e.g., ["Complete Blood Count"])
	TestTypeIDs        []int64                   `json:"test_type_ids"`        // Test type IDs (e.g., [1, 2, 3])
	SubCategoryIDs     []int64                   `json:"sub_category_ids"`     // Sub-category IDs (e.g., [1, 2])
	RequestedBy        string                    `json:"requested_by"`
	RequestedAt        string                    `json:"requested_at"`
}

// TechnoMedicPatientRequest represents patient information from TechnoMedic
type TechnoMedicPatientRequest struct {
	PatientID           string `json:"patient_id" validate:"required"`
	FullName            string `json:"full_name" validate:"required"`
	Sex                 string `json:"sex" validate:"required,oneof=M F"`
	Address             string `json:"address"`
	Birthdate           string `json:"birthdate" validate:"required"`
	MedicalRecordNumber string `json:"medical_record_number"`
	PhoneNumber         string `json:"phone_number"`
}

// TechnoMedicOrderResponse represents the response for order operations
type TechnoMedicOrderResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// TechnoMedicTestType represents test type information
type TechnoMedicTestType struct {
	ID           string `json:"id"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	Category     string `json:"category"`
	SubCategory  string `json:"sub_category"`
	SpecimenType string `json:"specimen_type"`
	Unit         string `json:"unit"`
}

// TechnoMedicSubCategory represents sub-category information
type TechnoMedicSubCategory struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description,omitempty"`
}

// TechnoMedicDoctor represents doctor information
type TechnoMedicDoctor struct {
	ID       int64  `json:"id"`
	Fullname string `json:"fullname"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

// TechnoMedicAnalyst represents analyst/analyzer information
type TechnoMedicAnalyst struct {
	ID       int64  `json:"id"`
	Fullname string `json:"fullname"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

// TechnoMedicGetOrderResponse represents the response for getting order details
type TechnoMedicGetOrderResponse struct {
	NoOrder          string                    `json:"no_order"`
	Status           string                    `json:"status"`
	Patient          TechnoMedicPatientRequest `json:"patient"`
	RequestedBy      string                    `json:"requested_by"`
	RequestedAt      string                    `json:"requested_at"`
	SubCategories    []SubCategory             `json:"sub_categories,omitempty"`
	ParametersResult []Results                 `json:"parameters_result,omitempty"` // For tests without sub-category
	CompletedAt      *string                   `json:"completed_at,omitempty"`
	VerifiedAt       *string                   `json:"verified_at,omitempty"`
	VerifiedBy       *string                   `json:"verified_by,omitempty"`
}

// SubCategory represents a test sub-category with results
type SubCategory struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	ParametersResult []Results `json:"parameters_result"`
}

// Results represents individual test results
type Results struct {
	ID           string `json:"id"`
	Code         string `json:"code"`
	CategoryName string `json:"category_name"`
	Value        string `json:"value"`
	SpecimenType string `json:"specimen_type"`
	Unit         string `json:"unit"`
	Ref          string `json:"ref"`
	Flag         string `json:"flag,omitempty"`
}
