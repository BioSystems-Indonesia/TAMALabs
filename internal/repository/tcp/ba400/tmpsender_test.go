package ba400

import (
	"context"
	"testing"
	"time"

	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestSendToBA400(t *testing.T) {
	t.Skip("need device")
	err := SendToBA400(context.Background(), []entity.Patient{
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
			Specimen: []entity.Specimen{
				{
					ID:             1,
					PatientID:      1,
					OrderID:        1,
					Type:           "SER",
					CollectionDate: time.Now().Format(time.RFC3339),
					ReceivedDate:   time.Time{},
					Source:         "",
					Condition:      "",
					Method:         "",
					Comments:       "",
					Barcode:        "123123",
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
					ObservationRequest: []entity.ObservationRequest{
						{
							ID:              1,
							SpecimenID:      1,
							TestCode:        "ACID GLYCO BIR",
							TestDescription: "ACID GLYCO BIR",
							CreatedAt:       time.Now(),
							UpdatedAt:       time.Now(),
						},
						{
							ID:              2,
							SpecimenID:      1,
							TestCode:        "ALBUMIN",
							TestDescription: "ALBUMIN",
							CreatedAt:       time.Now(),
							UpdatedAt:       time.Now(),
						},
					},
				},
			},
		},
		{
			ID:          2,
			FirstName:   "Dafa",
			LastName:    "Kedua",
			Birthdate:   time.Date(2002, time.October, 23, 0, 0, 0, 0, time.UTC),
			Sex:         "F",
			PhoneNumber: "",
			Location:    "",
			Address:     "",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Specimen: []entity.Specimen{
				{
					ID:             2,
					PatientID:      2,
					OrderID:        1,
					Type:           "SER",
					CollectionDate: time.Now().Format(time.RFC3339),
					ReceivedDate:   time.Time{},
					Source:         "",
					Condition:      "",
					Method:         "",
					Comments:       "",
					Barcode:        "456456",
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
					ObservationRequest: []entity.ObservationRequest{
						{
							ID:              3,
							SpecimenID:      2,
							TestCode:        "ACID GLYCO BIR",
							TestDescription: "ACID GLYCO BIR",
							CreatedAt:       time.Now(),
							UpdatedAt:       time.Now(),
						},
						{
							ID:              4,
							SpecimenID:      2,
							TestCode:        "ALBUMIN",
							TestDescription: "ALBUMIN",
							CreatedAt:       time.Now(),
							UpdatedAt:       time.Now(),
						},
						{
							ID:              5,
							SpecimenID:      2,
							TestCode:        "URIC ACID",
							TestDescription: "URIC ACID",
							CreatedAt:       time.Now(),
							UpdatedAt:       time.Now(),
						},
					},
				},
			},
		},
	}, entity.Device{
		ID:        1,
		Name:      "Device 1",
		IPAddress: "192.168.1.100",
		Port:      5678,
	}, false)

	if err != nil {
		t.Error(err)
	}
}

func Test_receiveResponse(t *testing.T) {
	type args struct {
		resp []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ORL_O34",
			args: args{
				resp: []byte(`MSH|^~\\&|BA200|Biosystems|Host|Host provider|20241223163505||ORL^O34^ORL_O34|ec3b41a9-77e3-4fd3-a2ca-f8f760dbda47|P|2.5.1|||ER|NE||UNICODE UTF-8|||LAB-28^IHE
MSA|AA|939b894f-a10a-4b35-9f82-95de095cc0c4`),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := receiveResponse(tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("receiveResponse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestByteLen(t *testing.T) {
	inp := []byte("MSH|^~\\&|LIS|Lab01|BA400|172.23.144.1|20250323111421||OML^O33^OML_O33|a5a25a9f-ce5d-40aa-9eed-b27870d84c2e|P|2.5.1|||ER|AL|ID|UNICODE UTF-8|||LAB-28^IHE\rPID|1|26|26||BANDI^SUBANDI||19640813170000|M\rSPM|1|20250321000043_01||SER^Serum^HL70369|||||||P\rORC|NW||||||||20250323111421\rTQ1|1||||||||R\rOBR|1|1||ACE^ACE^BA200|R\rORC|NW||||||||20250323111421\rTQ1|2||||||||R\rOBR|2|1||ALBUMIN^ALBUMIN^BA200|R\rORC|NW||||||||20250323111421\rTQ1|3||||||||R\rOBR|3|1||ALP-AMP^ALP-AMP^BA200|R\rORC|NW||||||||20250323111421\rTQ1|4||||||||R\rOBR|4|1||ALT-GPT^ALT-GPT^BA200|R\rORC|NW||||||||20250323111421\rTQ1|5||||||||R\rOBR|5|1||AMYLASE DIRECT^AMYLASE DIRECT^BA200|R\rORC|NW||||||||20250323111421\rTQ1|6||||||||R\rOBR|6|1||AST-GOT^AST-GOT^BA200|R\rORC|NW||||||||20250323111421\rTQ1|7||||||||R\rOBR|7|1||BASOPHILS^BASOPHILS^BA200|R\rORC|NW||||||||20250323111421\rTQ1|8||||||||R\rOBR|8|1||BILI DIRECT DPD^BILI DIRECT DPD^BA200|R\rORC|NW||||||||20250323111421\rTQ1|9||||||||R\rOBR|9|1||BILI TOTAL DPD^BILI TOTAL DPD^BA200|R\rORC|NW||||||||20250323111421\rTQ1|10||||||||R\rOBR|10|1||CALCIUM ARSENAZO^CALCIUM ARSENAZO^BA200|R\rORC|NW||||||||20250323111421\rTQ1|11||||||||R\rOBR|11|1||CHOLESTEROL^CHOLESTEROL^BA200|R\rORC|NW||||||||20250323111421\rTQ1|12||||||||R\rOBR|12|1||CHOLINESTERASE^CHOLINESTERASE^BA200|R\rORC|NW||||||||20250323111421\rTQ1|13||||||||R\rOBR|13|1||CK^CK^BA200|R\rORC|NW||||||||20250323111421\rTQ1|14||||||||R\rOBR|14|1||Diastole^Diastole^BA200|R\rORC|NW||||||||20250323111421\rTQ1|15||||||||R\rOBR|15|1||EOSINOPHILS^EOSINOPHILS^BA200|R\rORC|NW||||||||20250323111421\rTQ1|16||||||||R\rOBR|16|1||GAMMA-GT^GAMMA-GT^BA200|R\rORC|NW||||||||20250323111421\rTQ1|17||||||||R\rOBR|17|1||GLUCOSE^GLUCOSE^BA200|R\rORC|NW||||||||20250323111421\rTQ1|18||||||||R\rOBR|18|1||GRAN#^GRAN#^BA200|R\rORC|NW||||||||20250323111421\rTQ1|19||||||||R\rOBR|19|1||GRAN%^GRAN%^BA200|R\rORC|NW||||||||20250323111421\rTQ1|20||||||||R\rOBR|20|1||HBA1C-D%^HBA1C-D%^BA200|R\rORC|NW||||||||20250323111421\rTQ1|21||||||||R\rOBR|21|1||HBA1C-DIR^HBA1C-DIR^BA200|R\rORC|NW||||||||20250323111421\rTQ1|22||||||||R\rOBR|22|1||HCT^HCT^BA200|R\rORC|NW||||||||20250323111421\rTQ1|23||||||||R\rOBR|23|1||HDL DIRECT TOOS^HDL DIRECT TOOS^BA200|R\rORC|NW||||||||20250323111421\rTQ1|24||||||||R\rOBR|24|1||HEMATOCRIT^HEMATOCRIT^BA200|R\rORC|NW||||||||20250323111421\rTQ1|25||||||||R\rOBR|25|1||HEMOGLOBIN^HEMOGLOBIN^BA200|R\rORC|NW||||||||20250323111421\rTQ1|26||||||||R\rOBR|26|1||HGB^HGB^BA200|R\rORC|NW||||||||20250323111421\rTQ1|27||||||||R\rOBR|27|1||IRON FERROZINE^IRON FERROZINE^BA200|R\rORC|NW||||||||20250323111421\rTQ1|28||||||||R\rOBR|28|1||LCC^LCC^BA200|R\rORC|NW||||||||20250323111421\rTQ1|29||||||||R\rOBR|29|1||LCR^LCR^BA200|R\rORC|NW||||||||20250323111421\rTQ1|30||||||||R\rOBR|30|1||LDL DIRECT TOOS^LDL DIRECT TOOS^BA200|R\rORC|NW||||||||20250323111421\rTQ1|31||||||||R\rOBR|31|1||LYM#^LYM#^BA200|R\rORC|NW||||||||20250323111421\rTQ1|32||||||||R\rOBR|32|1||LYM%^LYM%^BA200|R\rORC|NW||||||||20250323111421\rTQ1|33||||||||R\rOBR|33|1||LYMPHOCYTES^LYMPHOCYTES^BA200|R\rORC|NW||||||||20250323111421\rTQ1|34||||||||R\rOBR|34|1||MAGNESIUM^MAGNESIUM^BA200|R\rORC|NW||||||||20250323111421\rTQ1|35||||||||R\rOBR|35|1||MCH^MCH^BA200|R\rORC|NW||||||||20250323111421\rTQ1|36||||||||R\rOBR|36|1||MCHC^MCHC^BA200|R\rORC|NW||||||||20250323111421\rTQ1|37||||||||R\rOBR|37|1||MCV^MCV^BA200|R\rORC|NW||||||||20250323111421\rTQ1|38||||||||R\rOBR|38|1||MID^MID^BA200|R\rORC|NW||||||||20250323111421\rTQ1|39||||||||R\rOBR|39|1||MID#^MID#^BA200|R\rORC|NW||||||||20250323111421\rTQ1|40||||||||R\rOBR|40|1||MONOCYTES^MONOCYTES^BA200|R\rORC|NW||||||||20250323111421\rTQ1|41||||||||R\rOBR|41|1||MPV^MPV^BA200|R\rORC|NW||||||||20250323111421\rTQ1|42||||||||R\rOBR|42|1||NEUTROPHILS^NEUTROPHILS^BA200|R\rORC|NW||||||||20250323111421\rTQ1|43||||||||R\rOBR|43|1||PCT^PCT^BA200|R\rORC|NW||||||||20250323111421\rTQ1|44||||||||R\rOBR|44|1||PDW^PDW^BA200|R\rORC|NW||||||||20250323111421\rTQ1|45||||||||R\rOBR|45|1||PLATELET COUNT^PLATELET COUNT^BA200|R\rORC|NW||||||||20250323111421\rTQ1|46||||||||R\rOBR|46|1||PLT^PLT^BA200|R\rORC|NW||||||||20250323111421\rTQ1|47||||||||R\rOBR|47|1||PROTEIN TOTALBIR^PROTEIN TOTALBIR^BA200|R\rORC|NW||||||||20250323111421\rTQ1|48||||||||R\rOBR|48|1||P_LCC^P_LCC^BA200|R\rORC|NW||||||||20250323111421\rTQ1|49||||||||R\rOBR|49|1||P_LCR^P_LCR^BA200|R\rORC|NW||||||||20250323111421\rTQ1|50||||||||R\rOBR|50|1||RBC^RBC^BA200|R\rORC|NW||||||||20250323111421\rTQ1|51||||||||R\rOBR|51|1||RDW^RDW^BA200|R\rORC|NW||||||||20250323111421\rTQ1|52||||||||R\rOBR|52|1||RDW_CV^RDW_CV^BA200|R\rORC|NW||||||||20250323111421\rTQ1|53||||||||R\rOBR|53|1||RDW_SD^RDW_SD^BA200|R\rORC|NW||||||||20250323111421\rTQ1|54||||||||R\rOBR|54|1||Sistole^Sistole^BA200|R\rORC|NW||||||||20250323111421\rTQ1|55||||||||R\rOBR|55|1||TRIGLYCERIDES^TRIGLYCERIDES^BA200|R\rORC|NW||||||||20250323111421\rTQ1|56||||||||R\rOBR|56|1||URIC ACID^URIC ACID^BA200|R\rORC|NW||||||||20250323111421\rTQ1|57||||||||R\rOBR|57|1||WBC^WBC^BA200|R\rSPM|2|20250321000043_02||SER^Serum^HL70369|||||||P\rORC|NW||||||||20250323111421\rTQ1|1||||||||R\rOBR|1|1||CREATININE^CREATININE^BA200|R\rORC|NW||||||||20250323111421\rTQ1|2||||||||R\rOBR|2|1||UREA-BUN-UV^UREA-BUN-UV^BA200|R\x1c\r\v")
	t.Log(string(inp))
	t.Log(len(inp))
	assert.True(t, false)
}
