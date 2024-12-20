package entity

import (
	"slices"
	"strings"
)

type Table struct {
	ID             string      `json:"id"`
	Name           string      `json:"name"`
	AdditionalInfo interface{} `json:"additional_info"`
}

type Tables []Table

func (t Tables) Find(id string) (Table, bool) {
	for _, table := range t {
		if table.ID == id {
			return table, true
		}
	}
	return Table{}, false
}

func (t Tables) FilterID(id []string) []Table {
	var tables []Table
	for _, table := range t {
		if slices.Contains(id, table.ID) {
			tables = append(tables, table)
		}
	}
	return tables
}

func (t Tables) FilterName(search string) []Table {
	var tables []Table
	for _, table := range t {
		if strings.Contains(strings.ToLower(table.Name), strings.ToLower(search)) {
			tables = append(tables, table)
		}
	}
	return tables
}

var TableSpecimenType = Tables{
	{
		ID: `SER`, Name: `Serum`,
	},
}

var TableSex = Tables{
	{
		ID: `M`, Name: `Male`,
	},
	{
		ID: `F`, Name: `Female`,
	},
	{
		ID: `U`, Name: `Unknown`,
	},
}

var TableSpecimenTest = Tables{
	{
		ID: `CHOLESTEROL`, Name: `CHOLESTEROL`,
	},
	{
		ID: `BUN`, Name: `BUN`,
	},
}

var TableWorkOrderStatus = Tables{
	{
		ID: string(WorkOrderStatusNew), Name: `NEW`,
	},
	{
		ID: string(WorkOrderCancelled), Name: `CANCELLED`,
	},
	{
		ID: string(WorkOrderStatusPending), Name: `PENDING`,
	},
	{
		ID: string(WorkOrderStatusCompleted), Name: `COMPLETED`,
	},
}

type ObservationInfo struct {
	Type ObservationTestType `json:"type"`
}

type ObservationTestType string

const (
	ObservationTestTypeStandard ObservationTestType = "STANDARD_TEST"
)

var TableObservationType = Tables{
	{ID: "ACE", Name: "ACE", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "ACID GLYCO BIR", Name: "ACID GLYCO BIR", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "ADA", Name: "ADA", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "ALBUMIN", Name: "ALBUMIN", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "ALBUMIN-MAU", Name: "ALBUMIN-MAU", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "ALP-AMP", Name: "ALP-AMP", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "ALP-DEA", Name: "ALP-DEA", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "ALPHA-GLUCOSIDAS", Name: "ALPHA-GLUCOSIDAS", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "ALT-GPT", Name: "ALT-GPT", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "AMMONIA", Name: "AMMONIA", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "AMYLASE DIRECT", Name: "AMYLASE DIRECT", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "AMYLASE EPS", Name: "AMYLASE EPS", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "AMYLASE PANCREAT", Name: "AMYLASE PANCREAT", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "APO AI", Name: "APO AI", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "APO B", Name: "APO B", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "ASO", Name: "ASO", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "AST-GOT", Name: "AST-GOT", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "AST-GOT-VERIF", Name: "AST-GOT-VERIF", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "ATIII", Name: "ATIII", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "B2 MICROGLOBULIN", Name: "B2 MICROGLOBULIN", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "B-HYDROXYBUTYRAT", Name: "B-HYDROXIBUTIRAT", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "BILI DIRECT DPD", Name: "BILI DIRECT DPD", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "BILI T NEWBORN", Name: "BILI T NEWBORN", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "BILI TOTAL DPD", Name: "BILI TOTAL DPD", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "BILIRUBIN DIRECT", Name: "BILIRUBIN DIRECT", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "BILIRUBIN TOTAL", Name: "BILIRUBIN TOTAL", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "CALCIUM ARSENAZO", Name: "CALCIUM ARSENAZO", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "CALCIUM CPC", Name: "CALCIUM CPC", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "CARBON DIOXIDE", Name: "CARBON DIOXIDE", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "CHOL HDL DIRECT", Name: "CHOL HDL DIRECT", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "CHOL LDL DIRECT", Name: "CHOL LDL DIRECT", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "CHOLESTEROL", Name: "CHOLESTEROL", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "CITRATE", Name: "CITRATE", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "CK", Name: "CK", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "CK-MB", Name: "CK-MB", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "COMPLEMENT C3BIR", Name: "COMPLEMENT C3BIR", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "COMPLEMENT C4BIR", Name: "COMPLEMENT C4BIR", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "CREATININE", Name: "CREATININE", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "CREATININE ENZ", Name: "CREATININE ENZ", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "CRP", Name: "CRP", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "CRPER", Name: "CRPER", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "CRPHS", Name: "CRPHS", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "D-DIMER", Name: "D-DIMER", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "DUMMY", Name: "DUMMY", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "ETHANOL", Name: "ETHANOL", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "FERRITIN", Name: "FERRITIN", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "FIBRINOGEN", Name: "FIBRINOGEN", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "FRUCTOSE", Name: "FRUCTOSE", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "G6PDH", Name: "G6PDH", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "GAMMA-GT", Name: "GAMMA-GT", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "GLUCOSE", Name: "GLUCOSE", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "GLUCOSE-HK", Name: "GLUCOSE-HK", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "GLUCOSE-VERIF", Name: "GLUCOSE-VERIF", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "HAPTOGLOBIN", Name: "HAPTOGLOBIN", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "HBA1C-DIRECT", Name: "HBA1C-DIR", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "HOMOCYSTEINE", Name: "HOMOCYSTEINE", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "IGA BIR", Name: "IGA BIR", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "IGG BIR", Name: "IGG BIR", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "IGM BIR", Name: "IGM BIR", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "IRON FERROZINE", Name: "IRON FERROZINE", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "LACTATE", Name: "LACTATE", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "LDH", Name: "LDH", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "LDH IFCC", Name: "LDH IFCC", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "LIPASE", Name: "LIPASE", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "LIPASE DGGR", Name: "LIPASE DGGR", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "MAGNESIUM", Name: "MAGNESIUM", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "NEFA", Name: "NEFA", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "OXALATE", Name: "OXALATE", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "PHOSPHORUS", Name: "PHOSPHORUS", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "PHOSPHORUS-VERIF", Name: "PHOSPHORUS-VERIF", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "PREALBUMIN BIR", Name: "PREALBUMIN BIR", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "PROTEIN TOTALBIR", Name: "PROTEIN TOTALBIR", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "PROTEIN URINE", Name: "PROTEIN URINE", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "RF", Name: "RF", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "TOTAL BILE ACIDS", Name: "TOTAL BILE ACIDS", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "TRANSFERRIN BIR", Name: "TRANSFERRIN BIR", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "TRIGLYCERIDES", Name: "TRIGLYCERIDES", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "UIBC", Name: "UIBC", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "UREA-BUN-UV", Name: "UREA-BUN-UV", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
	{ID: "URIC ACID", Name: "URIC ACID", AdditionalInfo: ObservationInfo{Type: ObservationTestTypeStandard}},
}

var TableList = map[string]Tables{
	"sex":               TableSex,
	"specimen-type":     TableSpecimenType,
	"specimen-test":     TableSpecimenTest,
	"observation-type":  TableObservationType,
	"work-order-status": TableWorkOrderStatus,
}
