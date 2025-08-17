package entity

import "time"

type KhanzaUpsertRequest struct {
	Insert []KhanzaResDT
	Update []KhanzaResDT
}

type DataTyp string

const (
	DataTypNumeric DataTyp = "NM"
	DataTypString  DataTyp = "ST"
	DataTypFT      DataTyp = "FT"
)

type Flag string

const (
	FlagLow      Flag = "L"
	FlagHigh     Flag = "H"
	FlagLowLow   Flag = "LL"
	FlagHighHigh Flag = "HH"
)

func NewKhanzaFlag(result TestResult) Flag {
	switch result.Abnormal {
	case HighResult:
		return FlagHigh
	case NormalResult:
		return ""
	case LowResult:
		return FlagLow
	case NoDataResult:
		return ""
	}

	return ""
}

// KhanzaResDT represents the KhanzaResDT table entity
type KhanzaResDT struct {
	ONO         string  `json:"ono"`
	OrderTestID string  `json:"order_testid"`
	TESTCD      string  `json:"test_cd"`
	TestNM      string  `json:"test_nm"`
	DataTyp     DataTyp `json:"data_typ"`
	ResultValue string  `json:"result_value"`
	ResultFT    string  `json:"result_ft"`
	Unit        string  `json:"unit"`
	Flag        Flag    `json:"flag"`
	RefRange    string  `json:"ref_range"`
}

// KhanzaLisOrder represents the lis_order table entity
type KhanzaLisOrder struct {
	ID           int64     `json:"id" db:"ID"`
	MessageDT    time.Time `json:"message_dt" db:"MESSAGE_DT"`
	OrderControl string    `json:"order_control" db:"ORDER_CONTROL"` // enum('NW','RP','CA')
	PID          string    `json:"pid" db:"PID"`
	PName        string    `json:"pname" db:"PNAME"`
	Address1     string    `json:"address1" db:"ADDRESS1"`
	Address2     string    `json:"address2" db:"ADDRESS2"`
	Address3     string    `json:"address3" db:"ADDRESS3"`
	Address4     string    `json:"address4" db:"ADDRESS4"`
	PType        string    `json:"ptype" db:"PTYPE"` // enum('IN','OP')
	BirthDT      time.Time `json:"birth_dt" db:"BIRTH_DT"`
	Sex          string    `json:"sex" db:"SEX"` // enum('1','2')
	ONO          string    `json:"ono" db:"ONO"`
	RequestDT    time.Time `json:"request_dt" db:"REQUEST_DT"`
	Source       string    `json:"source" db:"SOURCE"`
	Clinician    string    `json:"clinician" db:"CLINICIAN"`
	RoomNo       string    `json:"room_no" db:"ROOM_NO"`
	Priority     string    `json:"priority" db:"PRIORITY"` // enum('R','U')
	Comment      string    `json:"comment" db:"COMMENT"`
	VisitNo      string    `json:"visitno" db:"VISITNO"`
	OrderTestID  string    `json:"order_testid" db:"ORDER_TESTID"`
	Flag         string    `json:"flag" db:"FLAG"` // enum('0','1')
}

type KhanzaPatientSex string

const (
	KhanzaPatientSexMale   KhanzaPatientSex = "1"
	KhanzaPatientSexFemale KhanzaPatientSex = "2"
)

// KhanzaLabRequest represents the result of lab request query from main DB
type KhanzaLabRequest struct {
	NoOrder                string    `json:"noorder" db:"noorder"`
	NoRawat                string    `json:"no_rawat" db:"no_rawat"`
	NoRkmMedis             string    `json:"no_rkm_medis" db:"no_rkm_medis"`
	NmPasien               string    `json:"nm_pasien" db:"nm_pasien"`
	NmPerawatan            string    `json:"nm_perawatan" db:"nm_perawatan"`
	Pemeriksaan            string    `json:"pemeriksaan" db:"Pemeriksaan"`
	Satuan                 string    `json:"satuan" db:"satuan"`
	NilaiRujukanLD         string    `json:"nilai_rujukan_ld" db:"nilai_rujukan_ld"`
	KdPj                   string    `json:"kd_pj" db:"kd_pj"`
	NilaiRujukanLA         string    `json:"nilai_rujukan_la" db:"nilai_rujukan_la"`
	NilaiRujukanPD         string    `json:"nilai_rujukan_pd" db:"nilai_rujukan_pd"`
	NilaiRujukanPA         string    `json:"nilai_rujukan_pa" db:"nilai_rujukan_pa"`
	TglPermintaan          time.Time `json:"tgl_permintaan" db:"tgl_permintaan"`
	PngJawab               string    `json:"png_jawab" db:"png_jawab"`
	JamPermintaan          string    `json:"jam_permintaan" db:"jam_permintaan"`
	TglSampel              time.Time `json:"tgl_sampel" db:"tgl_sampel"`
	JamSampel              string    `json:"jam_sampel" db:"jam_sampel"`
	TglHasil               time.Time `json:"tgl_hasil" db:"tgl_hasil"`
	JamHasil               string    `json:"jam_hasil" db:"jam_hasil"`
	DokterPerujuk          string    `json:"dokter_perujuk" db:"dokter_perujuk"`
	NmDokter               string    `json:"nm_dokter" db:"nm_dokter"`
	NmPoli                 string    `json:"nm_poli" db:"nm_poli"`
	InformasiTambahan      string    `json:"informasi_tambahan" db:"informasi_tambahan"`
	DiagnosaKlinis         string    `json:"diagnosa_klinis" db:"diagnosa_klinis"`
}
