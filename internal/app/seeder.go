package app

import (
	"time"

	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

var seedPatient = []entity.Patient{
	{
		ID:          1,
		FirstName:   "Pasien",
		LastName:    "Pertama",
		Birthdate:   time.Date(1995, time.January, 1, 0, 0, 0, 0, time.UTC),
		Sex:         "M",
		PhoneNumber: "",
		Location:    "",
		Address:     "",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
	{
		ID:          2,
		FirstName:   "Pasien",
		LastName:    "Kedua",
		Birthdate:   time.Date(2002, time.October, 23, 0, 0, 0, 0, time.UTC),
		Sex:         "F",
		PhoneNumber: "",
		Location:    "",
		Address:     "",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
	{
		ID:          3,
		FirstName:   "Pasien",
		LastName:    "Ketiga",
		Birthdate:   time.Date(1998, time.February, 20, 0, 0, 0, 0, time.UTC),
		Sex:         "F",
		PhoneNumber: "",
		Location:    "",
		Address:     "",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
}

var seedDataTestType = []entity.TestType{
	// Hematology Tests
	// Red Series
	{Name: "RBC", Code: "RBC", Unit: "10^6/µL", LowRefRange: 4.5, HighRefRange: 5.9, Category: "Hematology", SubCategory: "Red Series", Description: "Red blood cell count.", Type: "Serum", Decimal: 2},
	{Name: "HEMOGLOBIN", Code: "HEMOGLOBIN", Unit: "g/dL", LowRefRange: 13.5, HighRefRange: 17.5, Category: "Hematology", SubCategory: "Red Series", Description: "Hemoglobin concentration in blood.", Type: "Serum", Decimal: 0},
	{Name: "HEMATOCRIT", Code: "HEMATOCRIT", Unit: "%", LowRefRange: 41, HighRefRange: 50, Category: "Hematology", SubCategory: "Red Series", Description: "Percentage of red blood cells in blood.", Type: "Serum", Decimal: 1},
	{Name: "MCV", Code: "MCV", Unit: "fL", LowRefRange: 80, HighRefRange: 100, Category: "Hematology", SubCategory: "Red Series", Description: "Mean corpuscular volume (size of red blood cells).", Type: "Serum", Decimal: 3},
	{Name: "MCH", Code: "MCH", Unit: "pg", LowRefRange: 27, HighRefRange: 31, Category: "Hematology", SubCategory: "Red Series", Description: "Mean corpuscular hemoglobin (average hemoglobin per red blood cell).", Type: "Serum", Decimal: 0},
	{Name: "MCHC", Code: "MCHC", Unit: "g/dL", LowRefRange: 32, HighRefRange: 36, Category: "Hematology", SubCategory: "Red Series", Description: "Mean corpuscular hemoglobin concentration.", Type: "Serum", Decimal: 4},
	{Name: "RDW", Code: "RDW", Unit: "%", LowRefRange: 11.5, HighRefRange: 14.5, Category: "Hematology", SubCategory: "Red Series", Description: "Red cell distribution width (variation in size of red blood cells).", Type: "Serum", Decimal: 5},
	{Name: "RDW_CV", Code: "RDW_CV", Unit: "%", LowRefRange: 11.5, HighRefRange: 14.5, Category: "Hematology", SubCategory: "Red Series", Description: "Coefficient of variation in red cell distribution width.", Type: "Serum", Decimal: 0},
	{Name: "RDW_SD", Code: "RDW_SD", Unit: "fL", LowRefRange: 37, HighRefRange: 54, Category: "Hematology", SubCategory: "Red Series", Description: "Standard deviation in red cell distribution width.", Type: "Serum", Decimal: 1},

	// White Series
	{Name: "WBC", Code: "WBC", Unit: "10^3/µL", LowRefRange: 4.0, HighRefRange: 11.0, Category: "Hematology", SubCategory: "White Series", Description: "White blood cell count.", Type: "Serum", Decimal: 2},
	{Name: "NEUTROPHILS", Code: "NEUTROPHILS", Unit: "%", LowRefRange: 40, HighRefRange: 70, Category: "Hematology", SubCategory: "White Series", Description: "Percentage of neutrophils in white blood cells.", Type: "Serum", Decimal: 0},
	{Name: "LYMPHOCYTES", Code: "LYMPHOCYTES", Unit: "%", LowRefRange: 20, HighRefRange: 40, Category: "Hematology", SubCategory: "White Series", Description: "Percentage of lymphocytes in white blood cells.", Type: "Serum", Decimal: 3},
	{Name: "MONOCYTES", Code: "MONOCYTES", Unit: "%", LowRefRange: 2, HighRefRange: 8, Category: "Hematology", SubCategory: "White Series", Description: "Percentage of monocytes in white blood cells.", Type: "Serum", Decimal: 4},
	{Name: "EOSINOPHILS", Code: "EOSINOPHILS", Unit: "%", LowRefRange: 1, HighRefRange: 4, Category: "Hematology", SubCategory: "White Series", Description: "Percentage of eosinophils in white blood cells.", Type: "Serum", Decimal: 0},
	{Name: "BASOPHILS", Code: "BASOPHILS", Unit: "%", LowRefRange: 0, HighRefRange: 1, Category: "Hematology", SubCategory: "White Series", Description: "Percentage of basophils in white blood cells.", Type: "Serum", Decimal: 5},
	{Name: "LCR", Code: "LCR", Unit: "%", LowRefRange: 0, HighRefRange: 10, Category: "Hematology", SubCategory: "White Series", Description: "Large cell ratio in white blood cells.", Type: "Serum", Decimal: 1},
	{Name: "LCC", Code: "LCC", Unit: "10^3/µL", LowRefRange: 0, HighRefRange: 0.4, Category: "Hematology", SubCategory: "White Series", Description: "Large cell count in white blood cells.", Type: "Serum", Decimal: 0},
	{Name: "MID", Code: "MID", Unit: "10^3/µL", LowRefRange: 0.1, HighRefRange: 0.9, Category: "Hematology", SubCategory: "White Series", Description: "Mid-size cells in blood (includes monocytes, eosinophils, and basophils).", Type: "Serum", Decimal: 2},
	{Name: "MID#", Code: "MID#", Unit: "10^3/µL", LowRefRange: 0.1, HighRefRange: 0.9, Category: "Hematology", SubCategory: "White Series", Description: "Absolute count of mid-size cells in blood.", Type: "Serum", Decimal: 3},

	// Platelets
	{Name: "PLATELET COUNT", Code: "PLATELET COUNT", Unit: "10^3/µL", LowRefRange: 150, HighRefRange: 450, Category: "Hematology", SubCategory: "Platelets", Description: "Platelet count in blood.", Type: "Serum", Decimal: 4},
	{Name: "MPV", Code: "MPV", Unit: "fL", LowRefRange: 7.5, HighRefRange: 11.5, Category: "Hematology", SubCategory: "Platelets", Description: "Mean platelet volume.", Type: "Serum", Decimal: 0},
	{Name: "PDW", Code: "PDW", Unit: "%", LowRefRange: 10, HighRefRange: 18, Category: "Hematology", SubCategory: "Platelets", Description: "Platelet distribution width (variation in size of platelets).", Type: "Serum", Decimal: 5},
	{Name: "PCT", Code: "PCT", Unit: "%", LowRefRange: 0.2, HighRefRange: 0.5, Category: "Hematology", SubCategory: "Platelets", Description: "Plateletcrit (total platelet mass in blood).", Type: "Serum", Decimal: 1},
	{Name: "P_LCR", Code: "P_LCR", Unit: "%", LowRefRange: 0, HighRefRange: 30, Category: "Hematology", SubCategory: "Platelets", Description: "Platelet large cell ratio.", Type: "Serum", Decimal: 0},
	{Name: "P_LCC", Code: "P_LCC", Unit: "10^3/µL", LowRefRange: 0, HighRefRange: 15, Category: "Hematology", SubCategory: "Platelets", Description: "Platelet large cell count.", Type: "Serum", Decimal: 2},

	// Observation Tests
	{Name: "ACE", Code: "ACE", Unit: "U/L", LowRefRange: 10, HighRefRange: 40, Category: "Observation", SubCategory: "Other", Description: "Angiotensin-converting enzyme activity.", Type: "Serum", Decimal: 3},
	{Name: "ACID GLYCO BIR", Code: "ACID GLYCO BIR", Unit: "mg/dL", LowRefRange: 5, HighRefRange: 15, Category: "Observation", SubCategory: "Other", Description: "Acid glycoprotein test.", Type: "Serum", Decimal: 0},
	{Name: "ADA", Code: "ADA", Unit: "U/L", LowRefRange: 10, HighRefRange: 45, Category: "Observation", SubCategory: "Other", Description: "Adenosine deaminase enzyme activity.", Type: "Serum", Decimal: 4},
	{Name: "ASO", Code: "ASO", Unit: "IU/mL", LowRefRange: 0, HighRefRange: 200, Category: "Observation", SubCategory: "Other", Description: "Anti-streptolysin O titer test.", Type: "Serum", Decimal: 5},
	{Name: "ATIII", Code: "ATIII", Unit: "%", LowRefRange: 70, HighRefRange: 120, Category: "Observation", SubCategory: "Other", Description: "Antithrombin III activity test.", Type: "Serum", Decimal: 1},
	{Name: "COMPLEMENT C3BIR", Code: "COMPLEMENT C3BIR", Unit: "mg/dL", LowRefRange: 90, HighRefRange: 180, Category: "Observation", SubCategory: "Other", Description: "Complement C3 test.", Type: "Serum", Decimal: 0},
	{Name: "COMPLEMENT C4BIR", Code: "COMPLEMENT C4BIR", Unit: "mg/dL", LowRefRange: 10, HighRefRange: 40, Category: "Observation", SubCategory: "Other", Description: "Complement C4 test.", Type: "Serum", Decimal: 2},
	{Name: "CRPER", Code: "CRPER", Unit: "", LowRefRange: 0, HighRefRange: 0, Category: "Observation", SubCategory: "Other", Description: "Observation for CRPER.", Type: "Serum", Decimal: 3},
	{Name: "DUMMY", Code: "DUMMY", Unit: "", LowRefRange: 0, HighRefRange: 0, Category: "Observation", SubCategory: "Other", Description: "Dummy observation for testing purposes.", Type: "Serum", Decimal: 0},

	// Biochemistry Tests
	{Name: "ALBUMIN", Code: "ALBUMIN", Unit: "g/dL", LowRefRange: 3.5, HighRefRange: 5.0, Category: "Biochemistry", SubCategory: "Liver Function", Description: "Measurement of serum albumin levels.", Type: "Serum", Decimal: 4},
	{Name: "ALBUMIN-MAU", Code: "ALBUMIN-MAU", Unit: "mg/g", LowRefRange: 0, HighRefRange: 30, Category: "Biochemistry", SubCategory: "Kidney Function", Description: "Microalbuminuria for kidney health assessment.", Type: "Urine", Decimal: 5},
	{Name: "ALP-AMP", Code: "ALP-AMP", Unit: "U/L", LowRefRange: 40, HighRefRange: 120, Category: "Biochemistry", SubCategory: "Liver Function", Description: "Alkaline phosphatase enzyme activity test (AMP method).", Type: "Serum", Decimal: 0},
	{Name: "ALP-DEA", Code: "ALP-DEA", Unit: "U/L", LowRefRange: 40, HighRefRange: 120, Category: "Biochemistry", SubCategory: "Liver Function", Description: "Alkaline phosphatase enzyme activity test (DEA method).", Type: "Serum", Decimal: 1},
	{Name: "ALPHA-GLUCOSIDAS", Code: "ALPHA-GLUCOSIDAS", Unit: "U/L", LowRefRange: 20, HighRefRange: 80, Category: "Biochemistry", SubCategory: "Metabolism", Description: "Alpha-glucosidase enzyme activity.", Type: "Serum", Decimal: 2},
	{Name: "ALT-GPT", Code: "ALT-GPT", Unit: "U/L", LowRefRange: 7, HighRefRange: 56, Category: "Biochemistry", SubCategory: "Liver Function", Description: "Alanine transaminase enzyme levels.", Type: "Serum", Decimal: 3},
	{Name: "AMMONIA", Code: "AMMONIA", Unit: "µg/dL", LowRefRange: 15, HighRefRange: 45, Category: "Biochemistry", SubCategory: "Other", Description: "Measurement of blood ammonia levels.", Type: "Serum", Decimal: 0},
	{Name: "AMYLASE DIRECT", Code: "AMYLASE DIRECT", Unit: "U/L", LowRefRange: 30, HighRefRange: 110, Category: "Biochemistry", SubCategory: "Other", Description: "Serum amylase enzyme activity test.", Type: "Serum", Decimal: 4},
	{Name: "AMYLASE EPS", Code: "AMYLASE EPS", Unit: "U/L", LowRefRange: 30, HighRefRange: 110, Category: "Biochemistry", SubCategory: "Other", Description: "Amylase test using EPS method.", Type: "Serum", Decimal: 5},
	{Name: "AMYLASE PANCREAT", Code: "AMYLASE PANCREAT", Unit: "U/L", LowRefRange: 30, HighRefRange: 110, Category: "Biochemistry", SubCategory: "Other", Description: "Pancreatic amylase test.", Type: "Serum", Decimal: 1},
	{Name: "APO AI", Code: "APO AI", Unit: "mg/dL", LowRefRange: 100, HighRefRange: 180, Category: "Biochemistry", SubCategory: "Other", Description: "Apolipoprotein AI levels.", Type: "Serum", Decimal: 0},
	{Name: "APO B", Code: "APO B", Unit: "mg/dL", LowRefRange: 60, HighRefRange: 120, Category: "Biochemistry", SubCategory: "Other", Description: "Apolipoprotein B levels.", Type: "Serum", Decimal: 2},
	{Name: "AST-GOT", Code: "AST-GOT", Unit: "U/L", LowRefRange: 10, HighRefRange: 40, Category: "Biochemistry", SubCategory: "Liver Function", Description: "Aspartate transaminase enzyme levels.", Type: "Serum", Decimal: 3},
	{Name: "AST-GOT-VERIF", Code: "AST-GOT-VERIF", Unit: "U/L", LowRefRange: 10, HighRefRange: 40, Category: "Biochemistry", SubCategory: "Liver Function", Description: "Aspartate transaminase levels verification.", Type: "Serum", Decimal: 4},
	{Name: "B2 MICROGLOBULIN", Code: "B2 MICROGLOBULIN", Unit: "mg/L", LowRefRange: 1.5, HighRefRange: 2.5, Category: "Biochemistry", SubCategory: "Other", Description: "Beta-2 microglobulin levels.", Type: "Serum", Decimal: 0},
	{Name: "B-HYDROXYBUTYRAT", Code: "B-HYDROXYBUTYRAT", Unit: "mmol/L", LowRefRange: 0.02, HighRefRange: 0.30, Category: "Biochemistry", SubCategory: "Metabolism", Description: "Beta-hydroxybutyrate levels.", Type: "Serum", Decimal: 5},
	{Name: "BILI DIRECT DPD", Code: "BILI DIRECT DPD", Unit: "mg/dL", LowRefRange: 0.0, HighRefRange: 0.4, Category: "Biochemistry", SubCategory: "Liver Function", Description: "Direct bilirubin test (DPD method).", Type: "Serum", Decimal: 1},
	{Name: "BILI T NEWBORN", Code: "BILI T NEWBORN", Unit: "mg/dL", LowRefRange: 0.0, HighRefRange: 12.0, Category: "Biochemistry", SubCategory: "Liver Function", Description: "Total bilirubin test for newborns.", Type: "Serum", Decimal: 0},
	{Name: "BILI TOTAL DPD", Code: "BILI TOTAL DPD", Unit: "mg/dL", LowRefRange: 0.2, HighRefRange: 1.2, Category: "Biochemistry", SubCategory: "Liver Function", Description: "Total bilirubin test (DPD method).", Type: "Serum", Decimal: 2},
	{Name: "BILIRUBIN DIRECT", Code: "BILIRUBIN DIRECT", Unit: "mg/dL", LowRefRange: 0.0, HighRefRange: 0.3, Category: "Biochemistry", SubCategory: "Liver Function", Description: "Measurement of direct bilirubin levels.", Type: "Serum", Decimal: 3},
	{Name: "BILIRUBIN TOTAL", Code: "BILIRUBIN TOTAL", Unit: "mg/dL", LowRefRange: 0.3, HighRefRange: 1.2, Category: "Biochemistry", SubCategory: "Liver Function", Description: "Measurement of total bilirubin levels.", Type: "Serum", Decimal: 4},
	{Name: "CALCIUM ARSENAZO", Code: "CALCIUM ARSENAZO", Unit: "mg/dL", LowRefRange: 8.5, HighRefRange: 10.5, Category: "Biochemistry", SubCategory: "Electrolytes", Description: "Calcium test using Arsenazo III method.", Type: "Serum", Decimal: 0},
	{Name: "CALCIUM CPC", Code: "CALCIUM CPC", Unit: "mg/dL", LowRefRange: 8.5, HighRefRange: 10.5, Category: "Biochemistry", SubCategory: "Electrolytes", Description: "Calcium test using CPC method.", Type: "Serum", Decimal: 5},
	{Name: "CARBON DIOXIDE", Code: "CARBON DIOXIDE", Unit: "mmol/L", LowRefRange: 23, HighRefRange: 29, Category: "Biochemistry", SubCategory: "Other", Description: "Measurement of serum bicarbonate levels.", Type: "Serum", Decimal: 1},
	{Name: "CHOL HDL DIRECT", Code: "CHOL HDL DIRECT", Unit: "mg/dL", LowRefRange: 40, HighRefRange: 60, Category: "Biochemistry", SubCategory: "Metabolism", Description: "High-density lipoprotein (HDL) cholesterol test.", Type: "Serum", Decimal: 0},
	{Name: "CHOL LDL DIRECT", Code: "CHOL LDL DIRECT", Unit: "mg/dL", LowRefRange: 70, HighRefRange: 130, Category: "Biochemistry", SubCategory: "Metabolism", Description: "Low-density lipoprotein (LDL) cholesterol test.", Type: "Serum", Decimal: 2},
	{Name: "CHOLESTEROL", Code: "CHOLESTEROL", Unit: "mg/dL", LowRefRange: 125, HighRefRange: 200, Category: "Biochemistry", SubCategory: "Metabolism", Description: "Total cholesterol levels.", Type: "Serum", Decimal: 3},
	{Name: "CITRATE", Code: "CITRATE", Unit: "mmol/L", LowRefRange: 0, HighRefRange: 0, Category: "Biochemistry", SubCategory: "Other", Description: "Citrate levels measurement.", Type: "Serum", Decimal: 4},
	{Name: "CK", Code: "CK", Unit: "U/L", LowRefRange: 38, HighRefRange: 174, Category: "Biochemistry", SubCategory: "Metabolism", Description: "Creatine kinase levels.", Type: "Serum", Decimal: 0},
	{Name: "CK-MB", Code: "CK-MB", Unit: "U/L", LowRefRange: 0, HighRefRange: 25, Category: "Biochemistry", SubCategory: "Metabolism", Description: "Creatine kinase-MB fraction levels.", Type: "Serum", Decimal: 5},
	{Name: "CRP", Code: "CRP", Unit: "mg/L", LowRefRange: 0, HighRefRange: 10, Category: "Biochemistry", SubCategory: "Other", Description: "C-reactive protein levels.", Type: "Serum", Decimal: 1},
	{Name: "CRPHS", Code: "CRPHS", Unit: "mg/L", LowRefRange: 0, HighRefRange: 10, Category: "Biochemistry", SubCategory: "Other", Description: "High-sensitivity C-reactive protein levels.", Type: "Serum", Decimal: 0},
	{Name: "D-DIMER", Code: "D-DIMER", Unit: "mg/L", LowRefRange: 0.0, HighRefRange: 0.5, Category: "Biochemistry", SubCategory: "Other", Description: "D-dimer test for clot formation.", Type: "Serum", Decimal: 2},
	{Name: "ETHANOL", Code: "ETHANOL", Unit: "mg/dL", LowRefRange: 0, HighRefRange: 50, Category: "Biochemistry", SubCategory: "Other", Description: "Ethanol levels in blood.", Type: "Serum", Decimal: 3},
	{Name: "FERRITIN", Code: "FERRITIN", Unit: "ng/mL", LowRefRange: 20, HighRefRange: 300, Category: "Biochemistry", SubCategory: "Other", Description: "Ferritin levels to assess iron stores.", Type: "Serum", Decimal: 4},
	{Name: "FIBRINOGEN", Code: "FIBRINOGEN", Unit: "mg/dL", LowRefRange: 200, HighRefRange: 400, Category: "Biochemistry", SubCategory: "Other", Description: "Fibrinogen levels in blood plasma.", Type: "Serum", Decimal: 0},
	{Name: "FRUCTOSE", Code: "FRUCTOSE", Unit: "mg/dL", LowRefRange: 0, HighRefRange: 0, Category: "Biochemistry", SubCategory: "Other", Description: "Fructose levels in the body.", Type: "Serum", Decimal: 5},
	{Name: "G6PDH", Code: "G6PDH", Unit: "U/L", LowRefRange: 4.6, HighRefRange: 13.5, Category: "Biochemistry", SubCategory: "Metabolism", Description: "G6PDH enzyme activity levels.", Type: "Serum", Decimal: 1},
	{Name: "GAMMA-GT", Code: "GAMMA-GT", Unit: "U/L", LowRefRange: 9, HighRefRange: 48, Category: "Biochemistry", SubCategory: "Liver Function", Description: "Gamma-glutamyl transferase enzyme levels.", Type: "Serum", Decimal: 0},
	{Name: "GLUCOSE", Code: "GLUCOSE", Unit: "mg/dL", LowRefRange: 70, HighRefRange: 99, Category: "Biochemistry", SubCategory: "Metabolism", Description: "Blood glucose levels.", Type: "Serum", Decimal: 2},
	{Name: "GLUCOSE-HK", Code: "GLUCOSE-HK", Unit: "mg/dL", LowRefRange: 70, HighRefRange: 99, Category: "Biochemistry", SubCategory: "Metabolism", Description: "Blood glucose using hexokinase method.", Type: "Serum", Decimal: 3},
	{Name: "GLUCOSE-VERIF", Code: "GLUCOSE-VERIF", Unit: "mg/dL", LowRefRange: 70, HighRefRange: 99, Category: "Biochemistry", SubCategory: "Metabolism", Description: "Blood glucose verification.", Type: "Serum", Decimal: 4},
	{Name: "HAPTOGLOBIN", Code: "HAPTOGLOBIN", Unit: "mg/dL", LowRefRange: 40, HighRefRange: 180, Category: "Biochemistry", SubCategory: "Other", Description: "Haptoglobin levels in blood.", Type: "Serum", Decimal: 0},
	{Name: "HBA1C-DIRECT", Code: "HBA1C-DIR", Unit: "%", LowRefRange: 4.0, HighRefRange: 5.6, Category: "Biochemistry", SubCategory: "Metabolism", Description: "Direct measurement of HbA1c for diabetes monitoring.", Type: "Serum", Decimal: 5},
	{Name: "HOMOCYSTEINE", Code: "HOMOCYSTEINE", Unit: "µmol/L", LowRefRange: 4, HighRefRange: 15, Category: "Biochemistry", SubCategory: "Metabolism", Description: "Homocysteine levels in blood.", Type: "Serum", Decimal: 1},
	{Name: "IGA BIR", Code: "IGA BIR", Unit: "mg/dL", LowRefRange: 70, HighRefRange: 400, Category: "Biochemistry", SubCategory: "Other", Description: "IgA antibody levels in blood.", Type: "Serum", Decimal: 0},
	{Name: "IGG BIR", Code: "IGG BIR", Unit: "mg/dL", LowRefRange: 700, HighRefRange: 1600, Category: "Biochemistry", SubCategory: "Other", Description: "IgG antibody levels in blood.", Type: "Serum", Decimal: 2},
	{Name: "IGM BIR", Code: "IGM BIR", Unit: "mg/dL", LowRefRange: 40, HighRefRange: 230, Category: "Biochemistry", SubCategory: "Other", Description: "IgM antibody levels in blood.", Type: "Serum", Decimal: 3},
	{Name: "IRON FERROZINE", Code: "IRON FERROZINE", Unit: "µg/dL", LowRefRange: 50, HighRefRange: 170, Category: "Biochemistry", SubCategory: "Metabolism", Description: "Iron levels using ferrozine method.", Type: "Serum", Decimal: 4},
	{Name: "LACTATE", Code: "LACTATE", Unit: "mmol/L", LowRefRange: 0.5, HighRefRange: 2.2, Category: "Biochemistry", SubCategory: "Metabolism", Description: "Blood lactate levels.", Type: "Serum", Decimal: 0},
	{Name: "LDH", Code: "LDH", Unit: "U/L", LowRefRange: 135, HighRefRange: 225, Category: "Biochemistry", SubCategory: "Other", Description: "Lactate dehydrogenase levels in blood.", Type: "Serum", Decimal: 5},
	{Name: "LDH IFCC", Code: "LDH IFCC", Unit: "U/L", LowRefRange: 135, HighRefRange: 225, Category: "Biochemistry", SubCategory: "Other", Description: "Lactate dehydrogenase using IFCC method.", Type: "Serum", Decimal: 1},
	{Name: "LIPASE", Code: "LIPASE", Unit: "U/L", LowRefRange: 0, HighRefRange: 160, Category: "Biochemistry", SubCategory: "Other", Description: "Lipase enzyme levels for pancreatic health.", Type: "Serum", Decimal: 0},
	{Name: "LIPASE DGGR", Code: "LIPASE DGGR", Unit: "U/L", LowRefRange: 0, HighRefRange: 160, Category: "Biochemistry", SubCategory: "Other", Description: "Lipase enzyme test using DGGR method.", Type: "Serum", Decimal: 2},
	{Name: "MAGNESIUM", Code: "MAGNESIUM", Unit: "mg/dL", LowRefRange: 1.7, HighRefRange: 2.2, Category: "Biochemistry", SubCategory: "Electrolytes", Description: "Magnesium levels in blood.", Type: "Serum", Decimal: 3},
	{Name: "NEFA", Code: "NEFA", Unit: "mmol/L", LowRefRange: 0, HighRefRange: 0, Category: "Biochemistry", SubCategory: "Metabolism", Description: "Non-esterified fatty acid levels.", Type: "Serum", Decimal: 4},
	{Name: "OXALATE", Code: "OXALATE", Unit: "mg/L", LowRefRange: 1, HighRefRange: 5, Category: "Biochemistry", SubCategory: "Other", Description: "Oxalate levels in blood or urine.", Type: "Serum", Decimal: 0},
	{Name: "PHOSPHORUS", Code: "PHOSPHORUS", Unit: "mg/dL", LowRefRange: 2.5, HighRefRange: 4.5, Category: "Biochemistry", SubCategory: "Other", Description: "Phosphorus levels in blood.", Type: "Serum", Decimal: 5},
	{Name: "PHOSPHORUS-VERIF", Code: "PHOSPHORUS-VERIF", Unit: "mg/dL", LowRefRange: 2.5, HighRefRange: 4.5, Category: "Biochemistry", SubCategory: "Other", Description: "Phosphorus levels verification.", Type: "Serum", Decimal: 1},
	{Name: "PREALBUMIN BIR", Code: "PREALBUMIN BIR", Unit: "mg/dL", LowRefRange: 15, HighRefRange: 36, Category: "Biochemistry", SubCategory: "Other", Description: "Prealbumin levels in blood.", Type: "Serum", Decimal: 0},
	{Name: "PROTEIN TOTALBIR", Code: "PROTEIN TOTALBIR", Unit: "g/dL", LowRefRange: 6.4, HighRefRange: 8.3, Category: "Biochemistry", SubCategory: "Other", Description: "Total protein levels in blood.", Type: "Serum", Decimal: 2},
	{Name: "PROTEIN URINE", Code: "PROTEIN URINE", Unit: "mg/dL", LowRefRange: 0, HighRefRange: 20, Category: "Biochemistry", SubCategory: "Other", Description: "Protein levels in urine.", Type: "Urine", Decimal: 3},
	{Name: "RF", Code: "RF", Unit: "IU/mL", LowRefRange: 0, HighRefRange: 14, Category: "Biochemistry", SubCategory: "Other", Description: "Rheumatoid factor levels in blood.", Type: "Serum", Decimal: 4},
	{Name: "TOTAL BILE ACIDS", Code: "TOTAL BILE ACIDS", Unit: "µmol/L", LowRefRange: 0, HighRefRange: 10, Category: "Biochemistry", SubCategory: "Liver Function", Description: "Total bile acid levels in blood.", Type: "Serum", Decimal: 0},
	{Name: "TRANSFERRIN BIR", Code: "TRANSFERRIN BIR", Unit: "mg/dL", LowRefRange: 200, HighRefRange: 400, Category: "Biochemistry", SubCategory: "Metabolism", Description: "Transferrin levels in blood.", Type: "Serum", Decimal: 5},
	{Name: "TRIGLYCERIDES", Code: "TRIGLYCERIDES", Unit: "mg/dL", LowRefRange: 50, HighRefRange: 150, Category: "Biochemistry", SubCategory: "Metabolism", Description: "Blood triglyceride levels.", Type: "Serum", Decimal: 1},
	{Name: "UIBC", Code: "UIBC", Unit: "µg/dL", LowRefRange: 155, HighRefRange: 355, Category: "Biochemistry", SubCategory: "Metabolism", Description: "Unsaturated iron-binding capacity in blood.", Type: "Serum", Decimal: 0},
	{Name: "UREA-BUN-UV", Code: "UREA-BUN-UV", Unit: "mg/dL", LowRefRange: 7, HighRefRange: 20, Category: "Biochemistry", SubCategory: "Kidney Function", Description: "Urea nitrogen levels in blood using UV method.", Type: "Serum", Decimal: 2},
	{Name: "URIC ACID", Code: "URIC ACID", Unit: "mg/dL", LowRefRange: 3.5, HighRefRange: 7.2, Category: "Biochemistry", SubCategory: "Metabolism", Description: "Uric acid levels in blood.", Type: "Serum", Decimal: 3},
}

var seedDevice = []entity.Device{
	{
		ID:        1,
		Name:      "Test Device",
		IPAddress: "192.168.1.100",
		Port:      5000,
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

	// Enzyme Activity Units
	{Value: "U/L", Base: Enzyme},

	// Immunology Units
	{Value: "IU/mL", Base: Immunology},
}
