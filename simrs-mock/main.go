package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/api/v1/emr/lab/list-new", handleGetLabList)
	http.HandleFunc("/api/v1/emr/lab/insert", handleInsertLabResult)
	http.HandleFunc("/api/v1/emr/lab/update-validasi", handleUpdateValidation)
	http.HandleFunc("/api/v1/emr/lab/insert-bulk", handleBulkLabInsert)

	log.Println("🏥 SIMRS Mock Server running on :4040")
	log.Fatal(http.ListenAndServe(":4040", nil))
}

type Metadata struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}
type LabListResponse struct {
	Response LabListData `json:"response"`
	Metadata Metadata    `json:"metadata"`
}

type LabListData struct {
	List []LabRegistration `json:"list"`
}

type LabRegistration struct {
	OrderDate       time.Time `json:"tgl"`
	LabNumber       int       `json:"no_lab"`
	MedicalRecordNo string    `json:"no_rm"`
	PatientName     string    `json:"nama"`
	BirthDate       time.Time `json:"tgl_lahir"`
	Gender          string    `json:"jenis_kelamin"`
	AgeDescription  string    `json:"umur"`
	Address         string    `json:"alamat"`
	Room            string    `json:"ruang"`
	Class           string    `json:"kelas"`
	InsuranceStatus string    `json:"status"`
	ReferringDoctor string    `json:"dokter_pengirim"`
	LabType         string    `json:"jenis_lab"`
	LISID           string    `json:"lis_id"`
	RoomID          string    `json:"id_ruangan"`
	RoomName        string    `json:"nama_ruangan"`
	InsuranceID     string    `json:"id_asuransi"`
	InsuranceName   string    `json:"nama_asuransi"`
	IsCITO          bool      `json:"cito,omitempty"`
	TestList        []LabTest `json:"list_test"`
}

type LabTest struct {
	DetailID    int          `json:"detail_id"`
	LabNumber   int          `json:"no_lab"`
	TestID      int          `json:"test_id"`
	TestName    string       `json:"nama_test"`
	LabType     string       `json:"jenis_lab"`
	TestType    string       `json:"jenis_test"`
	TestDetails []TestDetail `json:"detail_test"`
}

type TestDetail struct {
	PackageID int    `json:"paket_id"`
	Index     int    `json:"index"`
	Spacing   string `json:"spasi"`
	TestID    int    `json:"test_id"`
	TestName  string `json:"nama_test"`
}

func handleGetLabList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	resp := LabListResponse{
		Response: LabListData{
			List: []LabRegistration{
				{
					OrderDate:       time.Date(2023, 7, 31, 0, 12, 0, 0, time.UTC),
					LabNumber:       153310,
					MedicalRecordNo: "001063899",
					PatientName:     "Widad",
					BirthDate:       time.Date(1976, 5, 15, 0, 0, 0, 0, time.UTC),
					Gender:          "Laki-laki",
					AgeDescription:  "47 tahun, 2 bulan, 15 hari",
					Address:         "KP GEBANG RT 01/03 SANGIANG JAYA PERIUK",
					Room:            "",
					Class:           "",
					InsuranceStatus: "BPJS",
					ReferringDoctor: "dr. Arie Asnafi, Sp.U",
					LabType:         "pk",
					LISID:           "2307310001",
					RoomID:          "1",
					RoomName:        "Flamboyan",
					InsuranceID:     "2",
					InsuranceName:   "BPJS KESEHATAN",
					IsCITO:          false,
					TestList: []LabTest{
						{
							DetailID:  265836,
							LabNumber: 153304,
							TestID:    31,
							TestName:  "Darah Lengkap",
							LabType:   "pk",
							TestType:  "p", // p = paket
							TestDetails: []TestDetail{
								{
									PackageID: 31,
									Index:     1,
									Spacing:   "0",
									TestID:    16,
									TestName:  "Hemoglobin (HGB)",
								},
								{
									PackageID: 31,
									Index:     2,
									Spacing:   "0",
									TestID:    26,
									TestName:  "Leukocyte (WBC)",
								},
								{
									PackageID: 31,
									Index:     3,
									Spacing:   "0",
									TestID:    28,
									TestName:  "Erythrocyte (RBC)",
								},
								{
									PackageID: 31,
									Index:     4,
									Spacing:   "0",
									TestID:    29,
									TestName:  "Hematocrit (HCT)",
								},
								{
									PackageID: 31,
									Index:     5,
									Spacing:   "0",
									TestID:    30,
									TestName:  "Platelet (PLT)",
								},
								{
									PackageID: 31,
									Index:     6,
									Spacing:   "0",
									TestID:    32,
									TestName:  "MCV",
								},
								{
									PackageID: 31,
									Index:     7,
									Spacing:   "0",
									TestID:    33,
									TestName:  "MCH",
								},
								{
									PackageID: 31,
									Index:     8,
									Spacing:   "0",
									TestID:    34,
									TestName:  "MCHC",
								},
								{
									PackageID: 31,
									Index:     9,
									Spacing:   "0",
									TestID:    35,
									TestName:  "RDW",
								},
								{
									PackageID: 31,
									Index:     10,
									Spacing:   "0",
									TestID:    36,
									TestName:  "Neutrophil",
								},
								{
									PackageID: 31,
									Index:     11,
									Spacing:   "0",
									TestID:    37,
									TestName:  "Lymphocyte",
								},
								{
									PackageID: 31,
									Index:     12,
									Spacing:   "0",
									TestID:    38,
									TestName:  "Monocyte",
								},
								{
									PackageID: 31,
									Index:     13,
									Spacing:   "0",
									TestID:    39,
									TestName:  "Eosinophil",
								},
								{
									PackageID: 31,
									Index:     14,
									Spacing:   "0",
									TestID:    40,
									TestName:  "Basophil",
								},
							},
						},
					},
				},
			},
		},
		Metadata: Metadata{
			Message: "Ok",
			Code:    200,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type InsertLabResultRequest struct {
	LabNumber    int    `json:"no_lab"`
	TestID       int    `json:"test_id"`
	TestName     string `json:"nama_test"`
	ResultValue  string `json:"nilai"`
	Unit         string `json:"satuan"`
	NormalRange  string `json:"nilai_normal"`
	AbnormalFlag string `json:"flag"`
	ResultText   string `json:"keterangan"`
	PackageID    int    `json:"paket_id"`
	Index        int    `json:"index"`
	InsertedUser string `json:"inserted_user"`
	InsertedIP   string `json:"inserted_ip"`
}

func handleInsertLabResult(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read raw body
	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("❌ Error reading raw body: %v", err)
		http.Error(w, fmt.Sprintf("Error reading body: %v", err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("📝 Raw Request Body: %s", string(rawBody))

	// Decode JSON from raw body
	var reqBody InsertLabResultRequest
	if err := json.Unmarshal(rawBody, &reqBody); err != nil {
		log.Printf("❌ Error decoding request body: %v", err)
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	log.Printf("✅ Received lab result - Lab#%d, Test: %s (%d), Value: %s %s, Flag: %s, User: %s",
		reqBody.LabNumber,
		reqBody.TestName,
		reqBody.TestID,
		reqBody.ResultValue,
		reqBody.Unit,
		reqBody.AbnormalFlag,
		reqBody.InsertedUser,
	)

	fmt.Printf("📦 Parsed Struct: %+v\n", reqBody)

	resp := map[string]interface{}{
		"message": "Ok",
		"status":  200,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func handleBulkLabInsert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read raw body
	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("❌ Error reading raw body: %v", err)
		http.Error(w, fmt.Sprintf("Error reading body: %v", err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("📝 Raw Request Body: %s", string(rawBody))

	// // Decode JSON from raw body
	// var reqBody InsertLabResultRequest
	// if err := json.Unmarshal(rawBody, &reqBody); err != nil {
	// 	log.Printf("❌ Error decoding request body: %v", err)
	// 	http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
	// 	return
	// }

	// log.Printf("✅ Received lab result - Lab#%d, Test: %s (%d), Value: %s %s, Flag: %s, User: %s",
	// 	reqBody.LabNumber,
	// 	reqBody.TestName,
	// 	reqBody.TestID,
	// 	reqBody.ResultValue,
	// 	reqBody.Unit,
	// 	reqBody.AbnormalFlag,
	// 	reqBody.InsertedUser,
	// )

	// fmt.Printf("📦 Parsed Struct: %+v\n", reqBody)

	resp := map[string]interface{}{
		"message": "Ok",
		"status":  200,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func handleUpdateValidation(w http.ResponseWriter, r *http.Request) {
	resp := map[string]interface{}{
		"message": "Ok",
		"status":  200,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
