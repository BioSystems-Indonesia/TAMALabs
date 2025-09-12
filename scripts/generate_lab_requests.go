package main

import (
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"time"

	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var (
	firstNames = []string{
		"Ahmad", "Budi", "Siti", "Dewi", "Andi", "Rini", "Joko", "Sri", "Agus", "Lestari",
		"Wahyu", "Indira", "Bambang", "Ratna", "Hendro", "Wati", "Dedi", "Maya", "Rizki", "Sari",
		"Fajar", "Rina", "Yoga", "Dian", "Arif", "Lisa", "Bayu", "Novi", "Hendra", "Fitri",
		"Irwan", "Lia", "Yudi", "Mega", "Doni", "Tuti", "Eko", "Nita", "Reza", "Ani",
		"Gilang", "Putri", "Dimas", "Rara", "Fandi", "Sinta", "Galih", "Yuni", "Rio", "Vina",
	}

	lastNames = []string{
		"Pratama", "Wijaya", "Santoso", "Utomo", "Kusuma", "Sari", "Lestari", "Handayani",
		"Wibowo", "Kurniawan", "Rahayu", "Susanto", "Maharani", "Setiawan", "Purnama",
		"Andriani", "Nugroho", "Wardani", "Permata", "Saputra", "Indrayani", "Hartono",
		"Sartika", "Gunawan", "Puspita", "Hermawan", "Safitri", "Darmawan", "Anggraini",
		"Sudrajat", "Melati", "Wirawan", "Cahyani", "Saputri", "Sutrisno", "Kusumawati",
		"Budiman", "Fitriani", "Rahman", "Nuraini", "Hakim", "Suryani", "Hidayat",
		"Kartika", "Firmansyah", "Astuti", "Syahputra", "Dewanti", "Setiawati",
	}

	// Test types yang akan digunakan secara acak
	availableTestCodes = []string{
		"RBC", "HGB", "HCT", "MCV", "MCH", "MCHC", "RDW", "WBC", "NEUTROPHILS", "LYM%",
		"MONOCYTES", "EOSINOPHILS", "BASOPHILS", "PLT", "MPV", "PDW", "PCT",
		"ALBUMIN", "ALP-AMP", "ALT-GPT", "AST-GOT", "BILI DIRECT DPD", "BILI TOTAL DPD",
		"CALCIUM ARSENAZO", "CHOLESTEROL", "CK", "CREATININE", "CRP", "GLUCOSE",
		"UREA-BUN-UV", "URIC ACID", "TRIGLYCERIDES", "CHOL HDL DIRECT", "CHOL LDL DIRECT",
	}

	// Specimen types
	specimenTypes = []string{"SER", "URI", "CSF", "LIQ", "PLM"}

	// Cities in Indonesia
	cities = []string{
		"Jakarta", "Surabaya", "Bandung", "Medan", "Semarang", "Makassar", "Palembang",
		"Tangerang", "Depok", "Bekasi", "Bogor", "Batam", "Yogyakarta", "Malang",
		"Denpasar", "Bandar Lampung", "Balikpapan", "Samarinda", "Pontianak", "Manado",
	}
)

func main() {
	// Database path from the application
	const dbFileName = "./tmp/biosystem-lims.db"

	// Initialize database
	db, err := gorm.Open(sqlite.Open(dbFileName), &gorm.Config{})
	if err != nil {
		slog.Error("failed to connect database", "error", err)
		panic("failed to connect database")
	}

	slog.Info("Starting lab request generation...")

	// Generate 1000 lab requests
	if err := generateLabRequests(db, 1000); err != nil {
		slog.Error("failed to generate lab requests", "error", err)
		os.Exit(1)
	}

	slog.Info("Successfully generated 1000 lab requests!")
}

func generateLabRequests(db *gorm.DB, count int) error {
	// Get admin ID for created_by field
	var admin entity.Admin
	if err := db.First(&admin).Error; err != nil {
		return fmt.Errorf("failed to get admin: %w", err)
	}

	// Get available test types
	var testTypes []entity.TestType
	if err := db.Find(&testTypes).Error; err != nil {
		return fmt.Errorf("failed to get test types: %w", err)
	}

	if len(testTypes) == 0 {
		return fmt.Errorf("no test types available, please run seeder first")
	}

	// Create test type map for easy lookup
	testTypeMap := make(map[string]entity.TestType)
	for _, tt := range testTypes {
		testTypeMap[tt.Code] = tt
	}

	// Create a new random source
	randSource := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < count; i++ {
		if err := db.Transaction(func(tx *gorm.DB) error {
			// Generate patient
			patient := generatePatient(randSource)
			if err := tx.Create(&patient).Error; err != nil {
				return fmt.Errorf("failed to create patient %d: %w", i+1, err)
			}

			// Generate work order
			workOrder := generateWorkOrder(patient.ID, admin.ID, randSource)
			if err := tx.Create(&workOrder).Error; err != nil {
				return fmt.Errorf("failed to create work order %d: %w", i+1, err)
			}

			// Generate specimens (1-3 specimens per work order)
			specimenCount := randSource.Intn(3) + 1
			specimens := make([]entity.Specimen, 0, specimenCount)
			
			// Keep track of used specimen types to avoid UNIQUE constraint violation
			usedSpecimenTypes := make(map[string]bool)
			
			for j := 0; j < specimenCount && len(usedSpecimenTypes) < len(specimenTypes); j++ {
				// Get a specimen type that hasn't been used for this work order
				var specimenType string
				for {
					specimenType = specimenTypes[randSource.Intn(len(specimenTypes))]
					if !usedSpecimenTypes[specimenType] {
						usedSpecimenTypes[specimenType] = true
						break
					}
				}
				
				specimen := generateSpecimen(patient.ID, workOrder.ID, j+1, specimenType, randSource)
				specimens = append(specimens, specimen)
			}

			if err := tx.Create(&specimens).Error; err != nil {
				return fmt.Errorf("failed to create specimens for work order %d: %w", i+1, err)
			}

			// Generate observation requests for each specimen
			for _, specimen := range specimens {
				// Generate 1-5 tests per specimen
				testCount := randSource.Intn(5) + 1
				observationRequests := make([]entity.ObservationRequest, 0, testCount)

				// Use random test codes, avoiding duplicates for same specimen
				usedTestCodes := make(map[string]bool)
				for k := 0; k < testCount && len(usedTestCodes) < len(availableTestCodes); k++ {
					testCode := availableTestCodes[randSource.Intn(len(availableTestCodes))]
					if usedTestCodes[testCode] {
						continue
					}
					usedTestCodes[testCode] = true

					testType, exists := testTypeMap[testCode]
					if !exists {
						continue
					}

					observationRequest := entity.ObservationRequest{
						TestCode:        testCode,
						TestDescription: testType.Name,
						RequestedDate:   time.Now().Add(-time.Duration(randSource.Intn(30)) * 24 * time.Hour), // Random date in last 30 days
						ResultStatus:    "PENDING",
						SpecimenID:      int64(specimen.ID),
						CreatedAt:       time.Now(),
						UpdatedAt:       time.Now(),
					}

					observationRequests = append(observationRequests, observationRequest)
				}

				if len(observationRequests) > 0 {
					if err := tx.Create(&observationRequests).Error; err != nil {
						return fmt.Errorf("failed to create observation requests for specimen %d: %w", specimen.ID, err)
					}
				}
			}

			return nil
		}); err != nil {
			return err
		}

		// Log progress every 100 records
		if (i+1)%100 == 0 {
			slog.Info("Progress", "created", i+1, "total", count)
		}
	}

	return nil
}

func generatePatient(r *rand.Rand) entity.Patient {
	return entity.Patient{
		FirstName:   firstNames[r.Intn(len(firstNames))],
		LastName:    lastNames[r.Intn(len(lastNames))],
		Birthdate:   generateRandomBirthdate(r),
		Sex:         generateRandomSex(r),
		PhoneNumber: generatePhoneNumber(r),
		Location:    cities[r.Intn(len(cities))],
		Address:     generateAddress(r),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func generateWorkOrder(patientID int64, adminID int64, r *rand.Rand) entity.WorkOrder {
	// Generate random barcode
	barcode := fmt.Sprintf("WO%d%06d", time.Now().Unix()%10000, r.Intn(1000000))

	return entity.WorkOrder{
		Status:    entity.WorkOrderStatusNew,
		PatientID: patientID,
		CreatedBy: adminID,
		Barcode:   barcode,
		CreatedAt: time.Now().Add(-time.Duration(r.Intn(30)) * 24 * time.Hour), // Random date in last 30 days
		UpdatedAt: time.Now(),
	}
}

func generateSpecimen(patientID int64, workOrderID int64, sequenceNumber int, specimenType string, r *rand.Rand) entity.Specimen {
	barcode := fmt.Sprintf("SP%d%04d%02d", workOrderID, r.Intn(10000), sequenceNumber)

	collectionDate := time.Now().Add(-time.Duration(r.Intn(7)) * 24 * time.Hour)
	receivedDate := collectionDate.Add(time.Duration(r.Intn(24)) * time.Hour)

	return entity.Specimen{
		PatientID:      int(patientID),
		OrderID:        int(workOrderID),
		Type:           specimenType,
		CollectionDate: collectionDate.Format("2006-01-02 15:04:05"),
		ReceivedDate:   receivedDate,
		Source:         generateSpecimenSource(specimenType, r),
		Condition:      generateSpecimenCondition(r),
		Method:         "Standard Collection",
		Comments:       generateSpecimenComments(r),
		Barcode:        barcode,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

func generateRandomBirthdate(r *rand.Rand) time.Time {
	// Generate age between 1 and 90 years
	minAge := 1
	maxAge := 90
	age := r.Intn(maxAge-minAge) + minAge

	now := time.Now()
	birthYear := now.Year() - age

	// Random month and day
	month := r.Intn(12) + 1
	day := r.Intn(28) + 1 // Use 28 to avoid month boundary issues

	return time.Date(birthYear, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func generateRandomSex(r *rand.Rand) entity.PatientSex {
	sexes := []entity.PatientSex{entity.PatientSexMale, entity.PatientSexFemale}
	return sexes[r.Intn(len(sexes))]
}

func generatePhoneNumber(r *rand.Rand) string {
	// Generate Indonesian phone number format
	prefixes := []string{"081", "082", "083", "085", "087", "088", "089"}
	prefix := prefixes[r.Intn(len(prefixes))]

	// Generate 8 more digits
	suffix := ""
	for i := 0; i < 8; i++ {
		suffix += fmt.Sprintf("%d", r.Intn(10))
	}

	return prefix + suffix
}

func generateAddress(r *rand.Rand) string {
	streetNames := []string{
		"Jl. Sudirman", "Jl. Thamrin", "Jl. Gatot Subroto", "Jl. Kuningan", "Jl. Kemang",
		"Jl. Senopati", "Jl. Pramuka", "Jl. Veteran", "Jl. Diponegoro", "Jl. Ahmad Yani",
	}

	street := streetNames[r.Intn(len(streetNames))]
	number := r.Intn(999) + 1

	return fmt.Sprintf("%s No. %d", street, number)
}

func generateSpecimenSource(specimenType string, r *rand.Rand) string {
	sources := map[string][]string{
		"SER": {"Venous Blood", "Arterial Blood", "Capillary Blood"},
		"URI": {"Clean Catch", "Midstream", "Catheter"},
		"CSF": {"Lumbar Puncture", "Cisternal Puncture", "Ventricular Puncture"},
		"LIQ": {"Pleural Fluid", "Peritoneal Fluid", "Synovial Fluid"},
		"PLM": {"Plasma Sample", "EDTA Plasma", "Heparin Plasma"},
	}

	if sourcesForType, exists := sources[specimenType]; exists {
		return sourcesForType[r.Intn(len(sourcesForType))]
	}

	return "Standard"
}

func generateSpecimenCondition(r *rand.Rand) string {
	conditions := []string{
		"Good", "Acceptable", "Hemolyzed", "Lipemic", "Icteric", "Clotted", "Insufficient Volume",
	}

	// 70% chance of good condition
	if r.Float32() < 0.7 {
		return "Good"
	}

	return conditions[r.Intn(len(conditions))]
}

func generateSpecimenComments(r *rand.Rand) string {
	comments := []string{
		"Sample received in good condition",
		"No unusual observations",
		"Sample processed immediately",
		"Fasting sample as requested",
		"Post-prandial sample",
		"Sample collected after medication",
		"Emergency sample - stat processing",
		"Sample stored at appropriate temperature",
		"",
	}

	return comments[r.Intn(len(comments))]
}
