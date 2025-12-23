package app

import (
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	"golang.org/x/crypto/bcrypt"
)

var seedPatient = []entity.Patient{
	// {
	// 	ID:          1,
	// 	FirstName:   "Pasien",
	// 	LastName:    "Pertama",
	// 	Birthdate:   time.Date(1995, time.January, 1, 0, 0, 0, 0, time.UTC),
	// 	Sex:         "M",
	// 	PhoneNumber: "",
	// 	Location:    "",
	// 	Address:     "",
	// 	CreatedAt:   time.Now(),
	// 	UpdatedAt:   time.Now(),
	// },
	// {
	// 	ID:          2,
	// 	FirstName:   "Pasien",
	// 	LastName:    "Kedua",
	// 	Birthdate:   time.Date(2002, time.October, 23, 0, 0, 0, 0, time.UTC),
	// 	Sex:         "F",
	// 	PhoneNumber: "",
	// 	Location:    "",
	// 	Address:     "",
	// 	CreatedAt:   time.Now(),
	// 	UpdatedAt:   time.Now(),
	// },
	// {
	// 	ID:          3,
	// 	FirstName:   "Pasien",
	// 	LastName:    "Ketiga",
	// 	Birthdate:   time.Date(1998, time.February, 20, 0, 0, 0, 0, time.UTC),
	// 	Sex:         "F",
	// 	PhoneNumber: "",
	// 	Location:    "",
	// 	Address:     "",
	// 	CreatedAt:   time.Now(),
	// 	UpdatedAt:   time.Now(),
	// },
}

var seedDataTestType = []entity.TestType{
	// Hematology Tests - Assign to A15 device (ID: 1)
	{Name: "HEMOGLOBIN", Code: "HGB", Unit: "g/dL", LowRefRange: 12.0, HighRefRange: 17.0, Type: []entity.TestTypeSpecimenType{{Type: "WBL"}}, Decimal: 1, Category: "Hematology", SubCategory: "Red Series", Description: "Hemoglobin concentration in blood.", IsCalculatedTest: false, DeviceID: &[]int{1}[0]},
	{Name: "HEMATOCRIT", Code: "HCT", Unit: "%", LowRefRange: 33.0, HighRefRange: 45.0, Type: []entity.TestTypeSpecimenType{{Type: "WBL"}}, Decimal: 1, Category: "Hematology", SubCategory: "Red Series", Description: "Percentage of red blood cells in blood.", IsCalculatedTest: false, DeviceID: &[]int{1}[0]},
	{Name: "LEUKOSIT", Code: "WBC", Unit: "10^3/µL", LowRefRange: 4.00, HighRefRange: 10.00, Type: []entity.TestTypeSpecimenType{{Type: "WBL"}}, Decimal: 2, Category: "Hematology", SubCategory: "White Series", Description: "White blood cell count.", IsCalculatedTest: false, DeviceID: &[]int{1}[0]},
	{Name: "TROMBOSIT", Code: "PLT", Unit: "10^3/µL", LowRefRange: 130, HighRefRange: 400, Type: []entity.TestTypeSpecimenType{{Type: "WBL"}}, Decimal: 0, Category: "Hematology", SubCategory: "Platelets", Description: "Platelet count in blood.", IsCalculatedTest: false, DeviceID: &[]int{1}[0]},
	{Name: "ERITROSIT", Code: "RBC", Unit: "10^6/µL", LowRefRange: 4.50, HighRefRange: 5.50, Type: []entity.TestTypeSpecimenType{{Type: "WBL"}}, Decimal: 2, Category: "Hematology", SubCategory: "Red Series", Description: "Red blood cell count.", IsCalculatedTest: false, DeviceID: &[]int{1}[0]},
	{Name: "MCV", Code: "MCV", Unit: "µm3", LowRefRange: 75.0, HighRefRange: 100.0, Type: []entity.TestTypeSpecimenType{{Type: "WBL"}}, Decimal: 1, Category: "Hematology", SubCategory: "Red Series", Description: "Mean corpuscular volume.", IsCalculatedTest: false, DeviceID: &[]int{1}[0]},
	{Name: "MCH", Code: "MCH", Unit: "pg", LowRefRange: 25.0, HighRefRange: 35.0, Type: []entity.TestTypeSpecimenType{{Type: "WBL"}}, Decimal: 1, Category: "Hematology", SubCategory: "Red Series", Description: "Mean corpuscular hemoglobin.", IsCalculatedTest: false, DeviceID: &[]int{1}[0]},
	{Name: "MCHC", Code: "MCHC", Unit: "g/dL", LowRefRange: 31.0, HighRefRange: 38.0, Type: []entity.TestTypeSpecimenType{{Type: "WBL"}}, Decimal: 1, Category: "Hematology", SubCategory: "Red Series", Description: "Mean corpuscular hemoglobin concentration.", IsCalculatedTest: false, DeviceID: &[]int{1}[0]},
	{Name: "LYM%", Code: "LYM%", Unit: "%", LowRefRange: 15.00, HighRefRange: 50.00, Type: []entity.TestTypeSpecimenType{{Type: "WBL"}}, Decimal: 2, Category: "Hematology", SubCategory: "White Series", Description: "Percentage of lymphocytes in white blood cells.", IsCalculatedTest: false, DeviceID: &[]int{1}[0]},
	{Name: "GRA%", Code: "GRA%", Unit: "%", LowRefRange: 35.00, HighRefRange: 80.00, Type: []entity.TestTypeSpecimenType{{Type: "WBL"}}, Decimal: 2, Category: "Hematology", SubCategory: "White Series", Description: "Percentage of granulocytes in white blood cells.", IsCalculatedTest: false, DeviceID: &[]int{1}[0]},
	{Name: "MID%", Code: "MID%", Unit: "%", LowRefRange: 2.0, HighRefRange: 15.0, Type: []entity.TestTypeSpecimenType{{Type: "WBL"}}, Decimal: 1, Category: "Hematology", SubCategory: "White Series", Description: "Mid-size cells in blood.", IsCalculatedTest: false, DeviceID: &[]int{1}[0]},

	// Liver Function Tests (FUNGSI HATI) - No specific device assigned (general use)
	{Name: "SGOT/AST", Code: "AST", Unit: "U/L", LowRefRange: 0, HighRefRange: 42, Type: []entity.TestTypeSpecimenType{{Type: "SER"}}, Decimal: 0, Category: "Biochemistry", SubCategory: "Liver Function", Description: "Serum glutamic oxaloacetic transaminase / Aspartate aminotransferase.", IsCalculatedTest: false},
	{Name: "SGPT/ALT", Code: "ALT", Unit: "U/L", LowRefRange: 0, HighRefRange: 41, Type: []entity.TestTypeSpecimenType{{Type: "SER"}}, Decimal: 0, Category: "Biochemistry", SubCategory: "Liver Function", Description: "Serum glutamic pyruvic transaminase / Alanine aminotransferase.", IsCalculatedTest: false},
	{Name: "Total Protein", Code: "TP", Unit: "g/dL", LowRefRange: 6.60, HighRefRange: 8.80, Type: []entity.TestTypeSpecimenType{{Type: "SER"}}, Decimal: 2, Category: "Biochemistry", SubCategory: "Liver Function", Description: "Total protein in serum.", IsCalculatedTest: false},
	{Name: "Albumin", Code: "ALB", Unit: "g/dL", LowRefRange: 3.50, HighRefRange: 5.30, Type: []entity.TestTypeSpecimenType{{Type: "SER"}}, Decimal: 2, Category: "Biochemistry", SubCategory: "Liver Function", Description: "Albumin concentration in serum.", IsCalculatedTest: false},
	{Name: "Globulin", Code: "GLOB", Unit: "g/dL", LowRefRange: 0, HighRefRange: 0, Type: []entity.TestTypeSpecimenType{{Type: "SER"}}, Decimal: 1, Category: "Biochemistry", SubCategory: "Liver Function", Description: "Globulin concentration in serum.", IsCalculatedTest: true},

	// Kidney Function Tests (FUNGSI GINJAL) - No specific device assigned (general use)
	{Name: "Ureum", Code: "BUN", Unit: "mg/dL", LowRefRange: 10, HighRefRange: 50, Type: []entity.TestTypeSpecimenType{{Type: "SER"}}, Decimal: 0, Category: "Biochemistry", SubCategory: "Kidney Function", Description: "Blood urea nitrogen.", IsCalculatedTest: false},
	{Name: "Creatinin", Code: "CREA", Unit: "mg/dL", LowRefRange: 0.60, HighRefRange: 1.30, Type: []entity.TestTypeSpecimenType{{Type: "SER"}}, Decimal: 2, Category: "Biochemistry", SubCategory: "Kidney Function", Description: "Creatinine concentration in serum.", IsCalculatedTest: false},

	// Electrolytes (ELEKTROLIT) - No specific device assigned (general use)
	{Name: "Calcium Serum", Code: "CA", Unit: "mmol/L", LowRefRange: 2.10, HighRefRange: 2.70, Type: []entity.TestTypeSpecimenType{{Type: "SER"}}, Decimal: 2, Category: "Biochemistry", SubCategory: "Electrolytes", Description: "Calcium concentration in serum.", IsCalculatedTest: false},
	{Name: "Natrium Serum", Code: "NA", Unit: "mmol/L", LowRefRange: 136.0, HighRefRange: 145.0, Type: []entity.TestTypeSpecimenType{{Type: "SER"}}, Decimal: 1, Category: "Biochemistry", SubCategory: "Electrolytes", Description: "Sodium concentration in serum.", IsCalculatedTest: false},
	{Name: "Kalium Serum", Code: "K", Unit: "mmol/L", LowRefRange: 3.50, HighRefRange: 5.50, Type: []entity.TestTypeSpecimenType{{Type: "SER"}}, Decimal: 2, Category: "Biochemistry", SubCategory: "Electrolytes", Description: "Potassium concentration in serum.", IsCalculatedTest: false},
	{Name: "Chlorida Serum", Code: "CL", Unit: "mmol/L", LowRefRange: 96.0, HighRefRange: 108.0, Type: []entity.TestTypeSpecimenType{{Type: "SER"}}, Decimal: 1, Category: "Biochemistry", SubCategory: "Electrolytes", Description: "Chloride concentration in serum.", IsCalculatedTest: false},
}

var seedDevice = []entity.Device{
	{
		ID:          1,
		Name:        "A15",
		Type:        entity.DeviceTypeA15,
		ReceivePort: "10512",
		Path:        "C:\\Users\\Public\\Documents\\A15\\Import",
	},
	{
		ID:          2,
		Name:        "Swelab Alfa Plus",
		Type:        entity.DeviceTypeSwelabAlfa,
		ReceivePort: "4320",
	},
}

const (
	VolumeBase    = "Volume"
	MassBase      = "Mass"
	Percentage    = "Percentage"
	Concentration = "Concentration"
	Enzyme        = "Enzyme"
	Immunology    = "Immunology"
)

var seedUnits = []entity.Unit{
	// Volume-Based Units
	{Value: "10^6/µL", Base: VolumeBase},
	{Value: "10^3/µL", Base: VolumeBase},
	{Value: "pg", Base: VolumeBase},
	{Value: "fL", Base: VolumeBase},

	// Mass/Weight-Based Units
	{Value: "g/dL", Base: MassBase},
	{Value: "mg/dL", Base: MassBase},
	{Value: "µg/dL", Base: MassBase},
	{Value: "ng/mL", Base: MassBase},
	{Value: "mg/L", Base: MassBase},
	{Value: "mg/g", Base: MassBase},

	// Percentage Units
	{Value: "%", Base: Percentage},

	// Concentration Units
	{Value: "mmol/L", Base: Concentration},
	{Value: "µmol/L", Base: Concentration},
	{Value: "m mmol/L", Base: Concentration},

	// Enzyme Activity Units
	{Value: "U/L", Base: Enzyme},

	// Immunology Units
	{Value: "IU/mL", Base: Immunology},
}

var seedAdmin = []entity.Admin{
	initAdmin(),
}

func initAdmin() entity.Admin {
	defaultEmail := "admin@admin.com"
	const defaultUsername = "admin"
	const defaultPassword = "adminlishl7"

	hash, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	return entity.Admin{
		ID:           1,
		Username:     defaultUsername,
		Fullname:     "Admin",
		Email:        &defaultEmail,
		PasswordHash: string(hash),
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Roles: []entity.Role{
			{
				ID: 1,
			},
		},
	}
}

var seedRole = []entity.Role{
	{
		ID:          1,
		Name:        string(entity.RoleAdmin),
		Description: "Admin able to do all of LIMS features, including manage users. Only give admin permissions to highest authority user",
	},
	{
		ID:          2,
		Name:        string(entity.RoleDoctor),
		Description: "Doctor can be assigned as lab request doctor and able to approve result",
	},
	{
		ID:          3,
		Name:        string(entity.RoleAnalyst),
		Description: "Analyst can be assigned as lab request technician and able to perform lab request, but not approve result",
	},
}
