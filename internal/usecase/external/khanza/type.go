package khanzauc

type Request struct {
	Order Order `json:"order"`
}

type Order struct {
	Msh Msh `json:"msh"`
	PID PID `json:"pid"`
	OBR Obr `json:"obr"`
}

type Msh struct {
	Product string `json:"product"`
	Version string `json:"version"`
	UserID  string `json:"user_id"`
	Key     string `json:"key"`
}

type Obr struct {
	OrderControl  string   `json:"order_control"`
	Ptype         string   `json:"ptype"`
	RegNo         string   `json:"reg_no"`
	OrderLab      string   `json:"order_lab"`
	ProviderID    string   `json:"provider_id"`
	ProviderName  string   `json:"provider_name"`
	OrderDate     string   `json:"order_date"`
	ClinicianID   string   `json:"clinician_id"`
	ClinicianName string   `json:"clinician_name"`
	BangsalID     string   `json:"bangsal_id"`
	BangsalName   string   `json:"bangsal_name"`
	BedID         string   `json:"bed_id"`
	BedName       string   `json:"bed_name"`
	ClassID       string   `json:"class_id"`
	ClassName     string   `json:"class_name"`
	Cito          string   `json:"cito"`
	MedLegal      string   `json:"med_legal"`
	UserID        string   `json:"user_id"`
	Reserve1      string   `json:"reserve1"`
	Reserve2      string   `json:"reserve2"`
	Reserve3      string   `json:"reserve3"`
	Reserve4      string   `json:"reserve4"`
	OrderTest     []string `json:"order_test"`
}

type PID struct {
	Pmrn    string `json:"pmrn"`
	Pname   string `json:"pname"`
	Sex     string `json:"sex"`
	BirthDt string `json:"birth_dt"`
	Address string `json:"address"`
	NoTlp   string `json:"no_tlp"`
	NoHP    string `json:"no_hp"`
	Email   string `json:"email"`
	Nik     string `json:"nik"`
}

type Response struct {
	Result struct {
		OBX struct {
			OrderLab string `json:"order_lab"`
		} `json:"obx"`
	} `json:"result"`

	Response struct {
		Sample struct {
			ResultTest []ResponseResultTest `json:"result_test"`
		} `json:"sampel"` // "sampel" is required by the external API specification; this is intentional and not a typo.
	} `json:"response"`
}

type ResponseResultTest struct {
	TestID      string `json:"test_id"`
	NamaTest    string `json:"nama_test"`
	Hasil       string `json:"hasil"`
	NilaiNormal string `json:"nilai_normal"`
	Satuan      string `json:"satuan"`
	Flag        string `json:"flag"`
}
