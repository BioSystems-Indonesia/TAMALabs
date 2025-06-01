package app

import (
	"time"

	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"golang.org/x/crypto/bcrypt"
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
	{Name: "RBC", Code: "RBC", Unit: "10^6/µL", LowRefRange: 4.5, HighRefRange: 5.9, TypeDB: "SER", Decimal: 0, Category: "Hematology", SubCategory: "Red Series", Description: "Red blood cell count."},
	{Name: "HEMOGLOBIN", Code: "HGB", Unit: "g/dL", LowRefRange: 13.5, HighRefRange: 17.5, TypeDB: "", Decimal: 0, Category: "Hematology", SubCategory: "Red Series", Description: "Hemoglobin concentration in blood."},
	{Name: "HEMATOCRIT", Code: "HCT", Unit: "%", LowRefRange: 41, HighRefRange: 50, TypeDB: "", Decimal: 0, Category: "Hematology", SubCategory: "Red Series", Description: "Percentage of red blood cells in blood."},
	{Name: "MCV", Code: "MCV", Unit: "fL", LowRefRange: 80, HighRefRange: 100, TypeDB: "SER", Decimal: 0, Category: "Hematology", SubCategory: "Red Series", Description: "Mean corpuscular volume (size of red blood cells)."},
	{Name: "MCH", Code: "MCH", Unit: "pg", LowRefRange: 27, HighRefRange: 31, TypeDB: "SER", Decimal: 0, Category: "Hematology", SubCategory: "Red Series", Description: "Mean corpuscular hemoglobin (average hemoglobin per red blood cell)."},
	{Name: "MCHC", Code: "MCHC", Unit: "g/dL", LowRefRange: 32, HighRefRange: 36, TypeDB: "SER", Decimal: 0, Category: "Hematology", SubCategory: "Red Series", Description: "Mean corpuscular hemoglobin concentration."},
	{Name: "RDW", Code: "RDW", Unit: "%", LowRefRange: 11.5, HighRefRange: 14.5, TypeDB: "SER", Decimal: 0, Category: "Hematology", SubCategory: "Red Series", Description: "Red cell distribution width (variation in size of red blood cells)."},
	{Name: "RDW_CV", Code: "RDW_CV", Unit: "%", LowRefRange: 11.5, HighRefRange: 14.5, TypeDB: "SER", Decimal: 0, Category: "Hematology", SubCategory: "Red Series", Description: "Coefficient of variation in red cell distribution width."},
	{Name: "RDW_SD", Code: "RDW_SD", Unit: "fL", LowRefRange: 37, HighRefRange: 54, TypeDB: "SER", Decimal: 0, Category: "Hematology", SubCategory: "Red Series", Description: "Standard deviation in red cell distribution width."},
	{Name: "WBC", Code: "WBC", Unit: "10^3/µL", LowRefRange: 4, HighRefRange: 11, TypeDB: "SER", Decimal: 0, Category: "Hematology", SubCategory: "White Series", Description: "White blood cell count."},
	{Name: "NEUTROPHILS", Code: "NEUTROPHILS", Unit: "%", LowRefRange: 40, HighRefRange: 70, TypeDB: "SER", Decimal: 0, Category: "Hematology", SubCategory: "White Series", Description: "Percentage of neutrophils in white blood cells."},
	{Name: "LYM%", Code: "LYM%", Unit: "%", LowRefRange: 20, HighRefRange: 40, TypeDB: "", Decimal: 0, Category: "Hematology", SubCategory: "White Series", Description: "Percentage of lymphocytes in white blood cells."},
	{Name: "MONOCYTES", Code: "MONOCYTES", Unit: "%", LowRefRange: 2, HighRefRange: 8, TypeDB: "SER", Decimal: 0, Category: "Hematology", SubCategory: "White Series", Description: "Percentage of monocytes in white blood cells."},
	{Name: "EOSINOPHILS", Code: "EOSINOPHILS", Unit: "%", LowRefRange: 1, HighRefRange: 4, TypeDB: "SER", Decimal: 0, Category: "Hematology", SubCategory: "White Series", Description: "Percentage of eosinophils in white blood cells."},
	{Name: "BASOPHILS", Code: "BASOPHILS", Unit: "%", LowRefRange: 0, HighRefRange: 1, TypeDB: "SER", Decimal: 0, Category: "Hematology", SubCategory: "White Series", Description: "Percentage of basophils in white blood cells."},
	{Name: "LCR", Code: "LCR", Unit: "%", LowRefRange: 0, HighRefRange: 10, TypeDB: "SER", Decimal: 0, Category: "Hematology", SubCategory: "White Series", Description: "Large cell ratio in white blood cells."},
	{Name: "LCC", Code: "LCC", Unit: "10^3/µL", LowRefRange: 0, HighRefRange: 0.4, TypeDB: "SER", Decimal: 0, Category: "Hematology", SubCategory: "White Series", Description: "Large cell count in white blood cells."},
	{Name: "MID", Code: "MID", Unit: "%", LowRefRange: 1, HighRefRange: 15, TypeDB: "SER", Decimal: 0, Category: "Hematology", SubCategory: "White Series", Description: "Mid-size cells in blood (includes monocytes, eosinophils, and basophils)."},
	{Name: "MID#", Code: "MID#", Unit: "10^9/L", LowRefRange: 0.1, HighRefRange: 0.9, TypeDB: "SER", Decimal: 0, Category: "Hematology", SubCategory: "White Series", Description: "Absolute count of mid-size cells in blood."},
	{Name: "PLATELET COUNT", Code: "PLT", Unit: "10^3/µL", LowRefRange: 150, HighRefRange: 450, TypeDB: "", Decimal: 0, Category: "Hematology", SubCategory: "Platelets", Description: "Platelet count in blood."},
	{Name: "MPV", Code: "MPV", Unit: "fL", LowRefRange: 7.5, HighRefRange: 11.5, TypeDB: "SER", Decimal: 0, Category: "Hematology", SubCategory: "Platelets", Description: "Mean platelet volume."},
	{Name: "PDW", Code: "PDW", Unit: "%", LowRefRange: 10, HighRefRange: 18, TypeDB: "SER", Decimal: 0, Category: "Hematology", SubCategory: "Platelets", Description: "Platelet distribution width (variation in size of platelets)."},
	{Name: "PCT", Code: "PCT", Unit: "%", LowRefRange: 0.2, HighRefRange: 0.5, TypeDB: "SER", Decimal: 0, Category: "Hematology", SubCategory: "Platelets", Description: "Plateletcrit (total platelet mass in blood)."},
	{Name: "P_LCR", Code: "P_LCR", Unit: "%", LowRefRange: 0, HighRefRange: 30, TypeDB: "SER", Decimal: 0, Category: "Hematology", SubCategory: "Platelets", Description: "Platelet large cell ratio."},
	{Name: "P_LCC", Code: "P_LCC", Unit: "10^3/µL", LowRefRange: 0, HighRefRange: 15, TypeDB: "SER", Decimal: 0, Category: "Hematology", SubCategory: "Platelets", Description: "Platelet large cell count."},
	{Name: "ACE", Code: "ACE", Unit: "U/L", LowRefRange: 10, HighRefRange: 40, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Other", Description: "Angiotensin-converting enzyme activity."},
	{Name: "ACID GLYCO BIR", Code: "ACID GLYCO BIR", Unit: "mg/dL", LowRefRange: 5, HighRefRange: 15, TypeDB: "SER", Decimal: 0, Category: "Observation", SubCategory: "Other", Description: "Acid glycoprotein test."},
	{Name: "ADA", Code: "ADA", Unit: "U/L", LowRefRange: 10, HighRefRange: 45, TypeDB: "SER", Decimal: 0, Category: "Observation", SubCategory: "Other", Description: "Adenosine deaminase enzyme activity."},
	{Name: "ASO", Code: "ASO", Unit: "IU/mL", LowRefRange: 0, HighRefRange: 200, TypeDB: "SER", Decimal: 0, Category: "Observation", SubCategory: "Other", Description: "Anti-streptolysin O titer test."},
	{Name: "ATIII", Code: "ATIII", Unit: "%", LowRefRange: 70, HighRefRange: 120, TypeDB: "SER", Decimal: 0, Category: "Observation", SubCategory: "Other", Description: "Antithrombin III activity test."},
	{Name: "COMPLEMENT C3BIR", Code: "COMPLEMENT C3BIR", Unit: "mg/dL", LowRefRange: 90, HighRefRange: 180, TypeDB: "SER", Decimal: 0, Category: "Observation", SubCategory: "Other", Description: "Complement C3 test."},
	{Name: "COMPLEMENT C4BIR", Code: "COMPLEMENT C4BIR", Unit: "mg/dL", LowRefRange: 10, HighRefRange: 40, TypeDB: "SER", Decimal: 0, Category: "Observation", SubCategory: "Other", Description: "Complement C4 test."},
	{Name: "CRPER", Code: "CRPER", Unit: "", LowRefRange: 0, HighRefRange: 0, TypeDB: "SER", Decimal: 0, Category: "Observation", SubCategory: "Other", Description: "Observation for CRPER."},
	{Name: "DUMMY", Code: "DUMMY", Unit: "", LowRefRange: 0, HighRefRange: 0, TypeDB: "SER", Decimal: 0, Category: "Observation", SubCategory: "Other", Description: "Dummy observation for testing purposes."},
	{Name: "ALBUMIN", Code: "ALBUMIN", Unit: "g/L", LowRefRange: 3.5, HighRefRange: 5.2, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Liver Function", Description: "Measurement of serum albumin levels."},
	{Name: "ALBUMIN-MAU", Code: "ALBUMIN-MAU", Unit: "mg/g", LowRefRange: 0, HighRefRange: 30, TypeDB: "URI", Decimal: 0, Category: "Biochemistry", SubCategory: "Kidney Function", Description: "Microalbuminuria for kidney health assessment."},
	{Name: "ALP-AMP", Code: "ALP-AMP", Unit: "U/L", LowRefRange: 40, HighRefRange: 120, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Liver Function", Description: "Alkaline phosphatase enzyme activity test (AMP method)."},
	{Name: "ALP-DEA", Code: "ALP-DEA", Unit: "U/L", LowRefRange: 40, HighRefRange: 120, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Liver Function", Description: "Alkaline phosphatase enzyme activity test (DEA method)."},
	{Name: "ALPHA-GLUCOSIDAS", Code: "ALPHA-GLUCOSIDAS", Unit: "U/L", LowRefRange: 20, HighRefRange: 80, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Metabolism", Description: "Alpha-glucosidase enzyme activity."},
	{Name: "ALT-GPT", Code: "ALT-GPT", Unit: "U/L", LowRefRange: 7, HighRefRange: 56, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Liver Function", Description: "Alanine transaminase enzyme levels."},
	{Name: "AMMONIA", Code: "AMMONIA", Unit: "µg/dL", LowRefRange: 15, HighRefRange: 45, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Other", Description: "Measurement of blood ammonia levels."},
	{Name: "AMYLASE DIRECT", Code: "AMYLASE DIRECT", Unit: "U/L", LowRefRange: 30, HighRefRange: 110, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Other", Description: "Serum amylase enzyme activity test."},
	{Name: "AMYLASE EPS", Code: "AMYLASE EPS", Unit: "U/L", LowRefRange: 30, HighRefRange: 110, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Other", Description: "Amylase test using EPS method."},
	{Name: "AMYLASE PANCREAT", Code: "AMYLASE PANCREAT", Unit: "U/L", LowRefRange: 30, HighRefRange: 110, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Other", Description: "Pancreatic amylase test."},
	{Name: "APO AI", Code: "APO AI", Unit: "mg/dL", LowRefRange: 100, HighRefRange: 180, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Other", Description: "Apolipoprotein AI levels."},
	{Name: "APO B", Code: "APO B", Unit: "mg/dL", LowRefRange: 60, HighRefRange: 120, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Other", Description: "Apolipoprotein B levels."},
	{Name: "AST-GOT", Code: "AST-GOT", Unit: "U/L", LowRefRange: 10, HighRefRange: 40, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Liver Function", Description: "Aspartate transaminase enzyme levels."},
	{Name: "AST-GOT-VERIF", Code: "AST-GOT-VERIF", Unit: "U/L", LowRefRange: 10, HighRefRange: 40, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Liver Function", Description: "Aspartate transaminase levels verification."},
	{Name: "B2 MICROGLOBULIN", Code: "B2 MICROGLOBULIN", Unit: "mg/L", LowRefRange: 1.5, HighRefRange: 2.5, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Other", Description: "Beta-2 microglobulin levels."},
	{Name: "B-HYDROXYBUTYRAT", Code: "B-HYDROXYBUTYRAT", Unit: "mmol/L", LowRefRange: 0.02, HighRefRange: 0.3, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Metabolism", Description: "Beta-hydroxybutyrate levels."},
	{Name: "BILI DIRECT DPD", Code: "BILI DIRECT DPD", Unit: "mg/dL", LowRefRange: 0, HighRefRange: 0.4, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Liver Function", Description: "Direct bilirubin test (DPD method)."},
	{Name: "BILI T NEWBORN", Code: "BILI T NEWBORN", Unit: "mg/dL", LowRefRange: 0, HighRefRange: 12, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Liver Function", Description: "Total bilirubin test for newborns."},
	{Name: "BILI TOTAL DPD", Code: "BILI TOTAL DPD", Unit: "mg/dL", LowRefRange: 0.2, HighRefRange: 1.2, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Liver Function", Description: "Total bilirubin test (DPD method)."},
	{Name: "BILIRUBIN DIRECT", Code: "BILIRUBIN DIRECT", Unit: "mg/dL", LowRefRange: 0, HighRefRange: 0.3, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Liver Function", Description: "Measurement of direct bilirubin levels."},
	{Name: "BILIRUBIN TOTAL", Code: "BILIRUBIN TOTAL", Unit: "mg/dL", LowRefRange: 0.3, HighRefRange: 1.2, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Liver Function", Description: "Measurement of total bilirubin levels."},
	{Name: "CALCIUM ARSENAZO", Code: "CALCIUM ARSENAZO", Unit: "mg/dL", LowRefRange: 8.5, HighRefRange: 10.5, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Electrolytes", Description: "Calcium test using Arsenazo III method."},
	{Name: "CALCIUM CPC", Code: "CALCIUM CPC", Unit: "mg/dL", LowRefRange: 8.5, HighRefRange: 10.5, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Electrolytes", Description: "Calcium test using CPC method."},
	{Name: "CARBON DIOXIDE", Code: "CARBON DIOXIDE", Unit: "mmol/L", LowRefRange: 23, HighRefRange: 29, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Other", Description: "Measurement of serum bicarbonate levels."},
	{Name: "CHOL HDL DIRECT", Code: "CHOL HDL DIRECT", Unit: "mg/dL", LowRefRange: 40, HighRefRange: 60, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Metabolism", Description: "High-density lipoprotein (HDL) cholesterol test."},
	{Name: "CHOL LDL DIRECT", Code: "CHOL LDL DIRECT", Unit: "mg/dL", LowRefRange: 70, HighRefRange: 130, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Metabolism", Description: "Low-density lipoprotein (LDL) cholesterol test."},
	{Name: "CHOLESTEROL", Code: "CHOLESTEROL", Unit: "mg/dL", LowRefRange: 125, HighRefRange: 200, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Metabolism", Description: "Total cholesterol levels."},
	{Name: "CITRATE", Code: "CITRATE", Unit: "mmol/L", LowRefRange: 0, HighRefRange: 0, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Other", Description: "Citrate levels measurement."},
	{Name: "CK", Code: "CK", Unit: "U/L", LowRefRange: 38, HighRefRange: 174, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Metabolism", Description: "Creatine kinase levels."},
	{Name: "CK-MB", Code: "CK-MB", Unit: "U/L", LowRefRange: 0, HighRefRange: 25, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Metabolism", Description: "Creatine kinase-MB fraction levels."},
	{Name: "CRP", Code: "CRP", Unit: "mg/L", LowRefRange: 0, HighRefRange: 10, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Other", Description: "C-reactive protein levels."},
	{Name: "CRPHS", Code: "CRPHS", Unit: "mg/L", LowRefRange: 0, HighRefRange: 10, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Other", Description: "High-sensitivity C-reactive protein levels."},
	{Name: "D-DIMER", Code: "D-DIMER", Unit: "mg/L", LowRefRange: 0, HighRefRange: 0.5, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Other", Description: "D-dimer test for clot formation."},
	{Name: "ETHANOL", Code: "ETHANOL", Unit: "mg/dL", LowRefRange: 0, HighRefRange: 50, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Other", Description: "Ethanol levels in blood."},
	{Name: "FERRITIN", Code: "FERRITIN", Unit: "ng/mL", LowRefRange: 20, HighRefRange: 300, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Other", Description: "Ferritin levels to assess iron stores."},
	{Name: "FIBRINOGEN", Code: "FIBRINOGEN", Unit: "mg/dL", LowRefRange: 200, HighRefRange: 400, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Other", Description: "Fibrinogen levels in blood plasma."},
	{Name: "FRUCTOSE", Code: "FRUCTOSE", Unit: "mg/dL", LowRefRange: 0, HighRefRange: 0, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Other", Description: "Fructose levels in the body."},
	{Name: "G6PDH", Code: "G6PDH", Unit: "U/L", LowRefRange: 4.6, HighRefRange: 13.5, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Metabolism", Description: "G6PDH enzyme activity levels."},
	{Name: "GAMMA-GT", Code: "GAMMA-GT", Unit: "U/L", LowRefRange: 0, HighRefRange: 55, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Liver Function", Description: "Gamma-glutamyl transferase enzyme levels."},
	{Name: "GLUCOSE", Code: "GLUCOSE", Unit: "mg/dL", LowRefRange: 70, HighRefRange: 99, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Metabolism", Description: "Blood glucose levels."},
	{Name: "GLUCOSE-HK", Code: "GLUCOSE-HK", Unit: "mg/dL", LowRefRange: 70, HighRefRange: 99, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Metabolism", Description: "Blood glucose using hexokinase method."},
	{Name: "GLUCOSE-VERIF", Code: "GLUCOSE-VERIF", Unit: "mg/dL", LowRefRange: 70, HighRefRange: 99, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Metabolism", Description: "Blood glucose verification."},
	{Name: "HAPTOGLOBIN", Code: "HAPTOGLOBIN", Unit: "mg/dL", LowRefRange: 40, HighRefRange: 180, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Other", Description: "Haptoglobin levels in blood."},
	{Name: "HBA1C-DIRECT", Code: "HBA1C-DIR", Unit: "%", LowRefRange: 4, HighRefRange: 5.6, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Metabolism", Description: "Direct measurement of HbA1c for diabetes monitoring."},
	{Name: "HOMOCYSTEINE", Code: "HOMOCYSTEINE", Unit: "µmol/L", LowRefRange: 4, HighRefRange: 15, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Metabolism", Description: "Homocysteine levels in blood."},
	{Name: "IGA BIR", Code: "IGA BIR", Unit: "mg/dL", LowRefRange: 70, HighRefRange: 400, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Other", Description: "IgA antibody levels in blood."},
	{Name: "IGG BIR", Code: "IGG BIR", Unit: "mg/dL", LowRefRange: 700, HighRefRange: 1600, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Other", Description: "IgG antibody levels in blood."},
	{Name: "IGM BIR", Code: "IGM BIR", Unit: "mg/dL", LowRefRange: 40, HighRefRange: 230, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Other", Description: "IgM antibody levels in blood."},
	{Name: "IRON FERROZINE", Code: "IRON FERROZINE", Unit: "µg/dL", LowRefRange: 50, HighRefRange: 170, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Metabolism", Description: "Iron levels using ferrozine method."},
	{Name: "LACTATE", Code: "LACTATE", Unit: "mmol/L", LowRefRange: 0.5, HighRefRange: 2.2, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Metabolism", Description: "Blood lactate levels."},
	{Name: "LDH", Code: "LDH", Unit: "U/L", LowRefRange: 135, HighRefRange: 225, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Other", Description: "Lactate dehydrogenase levels in blood."},
	{Name: "LDH IFCC", Code: "LDH IFCC", Unit: "U/L", LowRefRange: 135, HighRefRange: 225, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Other", Description: "Lactate dehydrogenase using IFCC method."},
	{Name: "LIPASE", Code: "LIPASE", Unit: "U/L", LowRefRange: 0, HighRefRange: 160, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Other", Description: "Lipase enzyme levels for pancreatic health."},
	{Name: "LIPASE DGGR", Code: "LIPASE DGGR", Unit: "U/L", LowRefRange: 0, HighRefRange: 160, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Other", Description: "Lipase enzyme test using DGGR method."},
	{Name: "MAGNESIUM", Code: "MAGNESIUM", Unit: "mg/dL", LowRefRange: 1.7, HighRefRange: 2.2, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Electrolytes", Description: "Magnesium levels in blood."},
	{Name: "NEFA", Code: "NEFA", Unit: "mmol/L", LowRefRange: 0, HighRefRange: 0, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Metabolism", Description: "Non-esterified fatty acid levels."},
	{Name: "OXALATE", Code: "OXALATE", Unit: "mg/L", LowRefRange: 1, HighRefRange: 5, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Other", Description: "Oxalate levels in blood or urine."},
	{Name: "PHOSPHORUS", Code: "PHOSPHORUS", Unit: "mg/dL", LowRefRange: 2.5, HighRefRange: 4.5, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Other", Description: "Phosphorus levels in blood."},
	{Name: "PHOSPHORUS-VERIF", Code: "PHOSPHORUS-VERIF", Unit: "mg/dL", LowRefRange: 2.5, HighRefRange: 4.5, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Other", Description: "Phosphorus levels verification."},
	{Name: "PREALBUMIN BIR", Code: "PREALBUMIN BIR", Unit: "mg/dL", LowRefRange: 15, HighRefRange: 36, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Other", Description: "Prealbumin levels in blood."},
	{Name: "PROTEIN TOTALBIR", Code: "PROTEIN TOTALBIR", Unit: "g/dL", LowRefRange: 6.4, HighRefRange: 8.3, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Other", Description: "Total protein levels in blood."},
	{Name: "PROTEIN URINE", Code: "PROTEIN URINE", Unit: "mg/dL", LowRefRange: 0, HighRefRange: 20, TypeDB: "URI", Decimal: 0, Category: "Biochemistry", SubCategory: "Other", Description: "Protein levels in urine."},
	{Name: "RF", Code: "RF", Unit: "IU/mL", LowRefRange: 0, HighRefRange: 14, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Other", Description: "Rheumatoid factor levels in blood."},
	{Name: "TOTAL BILE ACIDS", Code: "TOTAL BILE ACIDS", Unit: "µmol/L", LowRefRange: 0, HighRefRange: 10, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Liver Function", Description: "Total bile acid levels in blood."},
	{Name: "TRANSFERRIN BIR", Code: "TRANSFERRIN BIR", Unit: "mg/dL", LowRefRange: 200, HighRefRange: 400, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Metabolism", Description: "Transferrin levels in blood."},
	{Name: "TRIGLYCERIDES", Code: "TRIGLYCERIDES", Unit: "mg/dL", LowRefRange: 50, HighRefRange: 150, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Metabolism", Description: "Blood triglyceride levels."},
	{Name: "UIBC", Code: "UIBC", Unit: "µg/dL", LowRefRange: 155, HighRefRange: 355, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Metabolism", Description: "Unsaturated iron-binding capacity in blood."},
	{Name: "UREA-BUN-UV", Code: "UREA-BUN-UV", Unit: "mg/dL", LowRefRange: 12.8, HighRefRange: 42.8, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Kidney Function", Description: "Urea nitrogen levels in blood using UV method."},
	{Name: "URIC ACID", Code: "URIC ACID", Unit: "mg/dL", LowRefRange: 3.5, HighRefRange: 7.2, TypeDB: "SER", Decimal: 0, Category: "Biochemistry", SubCategory: "Metabolism", Description: "Uric acid levels in blood."},
	{Name: "LDL DIRECT TOOS", Code: "LDL TOOS", Unit: "mg/dL", LowRefRange: 0, HighRefRange: 100, TypeDB: "", Decimal: 0, Category: "", SubCategory: "", Description: ""},
	{Name: "LDL DIRECT TOOS", Code: "LDL DIRECT TOOS", Unit: "mg/dL", LowRefRange: 0, HighRefRange: 100, TypeDB: "", Decimal: 0, Category: "Biochemistry", SubCategory: "", Description: ""},
	{Name: "HDL DIRECT TOOS", Code: "HDL DIRECT TOOS", Unit: "mg/dL", LowRefRange: 36, HighRefRange: 60, TypeDB: "", Decimal: 0, Category: "Biochemistry", SubCategory: "", Description: ""},
	{Name: "CREATININE", Code: "CREATININE", Unit: "mg/dL", LowRefRange: 0.5, HighRefRange: 1.2, TypeDB: "", Decimal: 0, Category: "Biochemistry", SubCategory: "", Description: ""},
	{Name: "HBA1C-DIR_NGSP", Code: "HBA1C-D%", Unit: "%", LowRefRange: 3.1, HighRefRange: 6.5, TypeDB: "", Decimal: 0, Category: "Biochemistry", SubCategory: "Biochemistry", Description: ""},
	{Name: "Sistole", Code: "Sistole", Unit: "mmHg", LowRefRange: 80, HighRefRange: 120, TypeDB: "", Decimal: 0, Category: "Blood Preasure", SubCategory: "", Description: ""},
	{Name: "Diastole", Code: "Diastole", Unit: "mmHg", LowRefRange: 60, HighRefRange: 80, TypeDB: "", Decimal: 0, Category: "Blood Preasure", SubCategory: "", Description: ""},
	{Name: "CHOLINESTERASE", Code: "CHOLINESTERASE", Unit: "U/L", LowRefRange: 3930, HighRefRange: 11500, TypeDB: "", Decimal: 0, Category: "Biochemistry", SubCategory: "Biochemistry", Description: ""},
	{Name: "LYM#", Code: "LYM#", Unit: "10^9/L", LowRefRange: 0.6, HighRefRange: 4.1, TypeDB: "", Decimal: 0, Category: "Hematology", SubCategory: "White Series", Description: ""},
	{Name: "GRAN%", Code: "GRAN%", Unit: "%", LowRefRange: 50, HighRefRange: 70, TypeDB: "", Decimal: 0, Category: "Hematology", SubCategory: "White Series", Description: ""},
	{Name: "GRAN#", Code: "GRAN#", Unit: "10^9/L", LowRefRange: 2, HighRefRange: 7.8, TypeDB: "", Decimal: 0, Category: "Hematology", SubCategory: "White Series", Description: ""},
	{Name: "HEMOGLOBIN", Code: "HEMOGLOBIN", Unit: "g/dL", LowRefRange: 13.5, HighRefRange: 17.5, TypeDB: "SER", Decimal: 0, Category: "Hematology", SubCategory: "Red Series", Description: "Hemoglobin concentration in blood."},
	{Name: "HEMATOCRIT", Code: "HEMATOCRIT", Unit: "%", LowRefRange: 41, HighRefRange: 50, TypeDB: "SER", Decimal: 0, Category: "Hematology", SubCategory: "Red Series", Description: "Percentage of red blood cells in blood."},
	{Name: "LYMPHOCYTES", Code: "LYMPHOCYTES", Unit: "%", LowRefRange: 20, HighRefRange: 40, TypeDB: "SER", Decimal: 0, Category: "Hematology", SubCategory: "White Series", Description: "Percentage of lymphocytes in white blood cells."},
	{Name: "PLATELET COUNT", Code: "PLATELET COUNT", Unit: "10^3/µL", LowRefRange: 150, HighRefRange: 450, TypeDB: "SER", Decimal: 0, Category: "Hematology", SubCategory: "Platelets", Description: "Platelet count in blood."},
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

var seedAdmin = []entity.Admin{
	initAdmin(),
}

func initAdmin() entity.Admin {
	const defaultEmail = "admin@admin.com"
	const defaultUsername = "admin"
	const defaultPassword = "123456"

	hash, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	return entity.Admin{
		ID:           1,
		Username:     defaultUsername,
		Fullname:     "First Admin",
		Email:        defaultEmail,
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
		Name:        string(entity.RoleAnalyzer),
		Description: "Analyzer can be assigned as lab request technician and able to perform lab request, but not approve result",
	},
}
