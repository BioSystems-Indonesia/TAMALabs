package nuha_simrs

import (
	"fmt"
	"strings"
	"time"
)

type NuhaTime struct {
	time.Time
}

// UnmarshalJSON accepts multiple date/time formats returned by Nuha SIMRS.
func (t *NuhaTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "null" || s == "" {
		// keep zero value
		t.Time = time.Time{}
		return nil
	}

	// Common layouts we expect from SIMRS: full RFC3339, space-separated, or date-only
	layouts := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006-01-02",
	}
	for _, l := range layouts {
		if tm, err := time.Parse(l, s); err == nil {
			t.Time = tm
			return nil
		}
	}

	return fmt.Errorf("unable to parse NuhaTime: %s", s)
}

// MarshalJSON outputs RFC3339 if value is non-zero, otherwise null
func (t NuhaTime) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", t.Time.Format(time.RFC3339))), nil
}

type LabListRequest struct {
	SessionID string `json:"session_id"`
	ValidFrom string `json:"valid_from"` // yyyy-mm-dd
	ValidTo   string `json:"valid_to"`   // yyyy-mm-dd
}

// InsertResultRequest represents request for inserting test result to Nuha SIMRS
type InsertResultRequest struct {
	SessionID      string `json:"session_id"`
	LabNumber      int    `json:"no_lab"`        // Lab number
	TestName       string `json:"nama_test"`     // Test name
	Result         string `json:"hasil"`         // Result value
	Unit           string `json:"satuan"`        // Unit
	ReferenceRange string `json:"nilai_rujukan"` // Reference range
	Abnormal       string `json:"abnormal"`      // Abnormal flag: "0"=Normal, "1"=Low, "2"=High
	Description    string `json:"keterangan"`    // Notes
	Notes          string `json:"catatan"`       // Additional notes
	TestID         int    `json:"test_id"`       // Test ID from Nuha SIMRS
	ResultText     string `json:"hasil_text"`    // Long text result
	PackageID      int    `json:"paket_id"`      // Package ID (0 if not from package)
	Spacing        string `json:"spasi"`         // Spacing
	Index          int    `json:"index"`         // Index order
	InsertedUser   string `json:"inserted_user"` // User who inserted
	InsertedIP     string `json:"inserted_ip"`   // IP address
}

// BatchInsertResultRequest represents batch request for inserting multiple test results
type BatchInsertResultRequest struct {
	SessionID string                  `json:"session_id"`
	Data      []BatchInsertResultItem `json:"data"`
}

// BatchInsertResultItem represents single test result item in batch request
type BatchInsertResultItem struct {
	LabNumber      int    `json:"no_lab"`
	TestName       string `json:"nama_test"`
	Result         string `json:"hasil"`
	ReferenceRange string `json:"nilai_rujukan"`
	Abnormal       string `json:"abnormal"`
	Unit           string `json:"satuan"`
	TestID         int    `json:"test_id"`
	PackageID      int    `json:"paket_id"` // 0 for non-package tests
	Index          int    `json:"index"`
	ResultText     string `json:"hasil_text"`
	InsertedUser   string `json:"inserted_user"`
	InsertedIP     string `json:"inserted_ip"`
}

// InsertResultResponse represents response from insert result API
type InsertResultResponse struct {
	Response InsertResultResponseData `json:"response"`
}

type InsertResultResponseData struct {
	Status  string `json:"status"`  // Success status
	Message string `json:"message"` // Response message
}

type LabListResponse struct {
	Response LabListData `json:"response"`
	Metadata Metadata    `json:"metadata"`
}

type LabListData struct {
	List []LabRegistration `json:"list"`
}

type LabRegistration struct {
	OrderDate       NuhaTime  `json:"tgl"`
	LabNumber       int       `json:"no_lab"`
	MedicalRecordNo string    `json:"no_rm"`
	PatientName     string    `json:"nama"`
	BirthDate       NuhaTime  `json:"tgl_lahir"`
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

type InsertLabResultRequest struct {
	SessionID    string `json:"session_id"`
	LabNumber    int    `json:"no_lab"`
	TestName     string `json:"nama_test"`
	ResultValue  string `json:"hasil"`
	Unit         string `json:"satuan"`
	Reference    string `json:"nilai_rujukan"`
	AbnormalFlag string `json:"abnormal"` // 0,1,2,3
	Note         string `json:"keterangan"`
	Comment      string `json:"catatan"`
	TestID       int    `json:"test_id"`
	PackageID    int    `json:"paket_id"`
	Spacing      string `json:"spasi"`
	Index        int    `json:"index"`
	ResultText   string `json:"hasil_text"`
	InsertedUser string `json:"inserted_user"`
	InsertedIP   string `json:"inserted_ip"`
}

type SimpleResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type LabValidationRequest struct {
	SessionID        string `json:"session_id"`
	LabNumber        int    `json:"no_lab"`
	ValidationStatus string `json:"status_validasi"`
	ValidatedBy      string `json:"user_validasi"`
}

type LabDetailRequest struct {
	SessionID string `json:"session_id"`
	LabNumber int    `json:"no_lab"`
}

type LabDetailResponse struct {
	Response LabRegistration `json:"response"`
	Metadata Metadata        `json:"metadata"`
}

type UpdateLabFlagRequest struct {
	SessionID    string `json:"session_id"`
	LabNumber    int    `json:"no_lab"`
	Flag         int    `json:"flag"`
	InsertedUser string `json:"inserted_user"`
	InsertedIP   string `json:"inserted_ip"`
}

type UpdateLabFlagResponse struct {
	Response LabFlagLog `json:"response"`
	Metadata Metadata   `json:"metadata"`
}

type LabFlagLog struct {
	LogID        int      `json:"id_log"`
	LabNumber    int      `json:"no_lab"`
	FetchedAt    NuhaTime `json:"fetched_at"`
	Flag         bool     `json:"flag"`
	InsertedUser string   `json:"inserted_user"`
	InsertedDate NuhaTime `json:"inserted_date"`
	InsertedIP   string   `json:"inserted_ip"`
	UpdatedUser  string   `json:"updated_user"`
	UpdatedDate  NuhaTime `json:"updated_date"`
	UpdatedIP    string   `json:"updated_ip"`
}

type InsertBulkResultRequest struct {
	SessionID string              `json:"session_id"`
	Data      []BulkLabResultItem `json:"data"`
}

type BulkLabResultItem struct {
	LabNumber    int    `json:"no_lab"`
	TestName     string `json:"nama_test"`
	ResultValue  string `json:"hasil"`
	Reference    string `json:"nilai_rujukan"`
	AbnormalFlag string `json:"abnormal"`
	Unit         string `json:"satuan"`
	TestID       int    `json:"test_id"`
	PackageID    *int   `json:"paket_id"`
	Index        int    `json:"index"`
	ResultText   string `json:"hasil_text"`
	InsertedUser string `json:"inserted_user"`
	InsertedIP   string `json:"inserted_ip"`
}

type InsertBulkResultResponse struct {
	Response []LabResultRecord `json:"response"`
	Metadata Metadata          `json:"metadata"`
}

type LabResultRecord struct {
	ResultID       int      `json:"hasil_id"`
	LabNumber      int      `json:"no_lab"`
	TestName       string   `json:"nama_test"`
	ResultValue    string   `json:"hasil"`
	Unit           string   `json:"satuan"`
	Reference      string   `json:"nilai_rujukan"`
	AbnormalFlag   string   `json:"abnormal"`
	Note           string   `json:"keterangan"`
	Comment        string   `json:"catatan"`
	TestID         int      `json:"test_id"`
	PackageID      int      `json:"paket_id"`
	Index          int      `json:"index"`
	ResultText     string   `json:"hasil_text"`
	Spacing        string   `json:"spasi"`
	Status         string   `json:"status"`
	ValidationFlag string   `json:"validasi"`
	ValidatedBy    string   `json:"user_validasi"`
	ValidatedAt    NuhaTime `json:"tanggal_validasi"`
	InsertedUser   string   `json:"inserted_user"`
	InsertedDate   NuhaTime `json:"inserted_date"`
	InsertedIP     string   `json:"inserted_ip"`
	UpdatedUser    string   `json:"updated_user"`
	UpdatedDate    NuhaTime `json:"updated_date"`
	UpdatedIP      string   `json:"updated_ip"`
}

type Metadata struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}
