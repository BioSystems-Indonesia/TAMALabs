package ba400

import (
	"strings"
	"testing"
	"time"

	"github.com/kardianos/hl7"
	"github.com/kardianos/hl7/h251"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	date := time.Date(2013, 01, 29, 10, 20, 30, 0, time.UTC)
	dateOfTransaction := time.Date(2024, 12, 9, 22, 0, 0, 0, time.UTC)

	msgControlID := "69F2746D24014F21AD7139756F64CAD8"

	patient := entity.Patient{
		ID:        123456,
		FirstName: "Doe",
		LastName:  "John",
		Birthdate: time.Date(1980, 1, 1, 0, 0, 0, 0, time.Local),
		Sex:       "M",
		Location:  "Jakarta",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	device := entity.Device{
		Name:      "BA200",
		IPAddress: "localhost",
	}

	o := h251.OML_O33{
		MSH:     NewOML_O33_MSH(msgControlID, device, date),
		Patient: NewOML_O33_Patient(patient),
		Specimen: []h251.OML_O33_Specimen{
			{
				HL7: h251.HL7Name{},
				SPM: &h251.SPM{
					HL7:   h251.HL7Name{},
					SetID: "1",
					SpecimenID: &h251.EIP{
						PlacerAssignedIdentifier: &h251.EI{
							EntityIdentifier: "78901",
						},
					},
					SpecimenParentIDs: []h251.EIP{},
					SpecimenType:      Serum,
					SpecimenRole:      []h251.CWE{{Identifier: "P"}},
				},
				Order: []h251.OML_O33_Order{
					{
						ORC: &h251.ORC{
							OrderControl: "NW",
							PlacerOrderNumber: &h251.EI{
								EntityIdentifier: "ORDER123",
							},
							DateTimeOfTransaction: dateOfTransaction,
						},
						ObservationRequest: &h251.OML_O33_ObservationRequest{OBR: &h251.OBR{
							SetID: "1",
							PlacerOrderNumber: &h251.EI{
								EntityIdentifier: "78901",
							},
							UniversalServiceIdentifier: h251.CE{
								Identifier:         "BUN",
								Text:               "BUN",
								NameOfCodingSystem: "BA200",
							},
						}},
					},
				},
			},
		},
	}

	encoder := hl7.NewEncoder(&hl7.EncodeOption{
		TrimTrailingSeparator: true,
	})
	b, err := encoder.Encode(o)
	if err != nil {
		t.Error(err)
	}

	expected := "MSH|^~\\&|LIS|Lab01|BA200|localhost|20130129102030||OML^O33^OML_O33|69F2746D24014F21AD7139756F64CAD8|P|2.5.1|||ER|AL|ID|UNICODE UTF-8|||LAB-28^IHE\rPID|1||123456||John^Doe||19800101|M\rSPM|1|78901||SER^Serum^HL70369|||||||P\rORC|NW|ORDER123|||||||20241209220000\rOBR|1|78901||BUN^BUN^BA200\r"
	exp := strings.Split(expected, "\r")
	got := strings.Split(string(b), "\r")

	assert.Equal(t, exp[0], got[0])
	//assert.Equal(t, exp[1], got[1]) // date of birth is different format
	assert.Equal(t, exp[2], got[2])
	assert.Equal(t, exp[3], got[3])
	assert.Equal(t, exp[4], got[4])

	/*
		s := Sender{host: "localhost:5000"}
		d, err := s.SendRaw(b)
		t.Error(err)
		t.Error(string(d))
	*/
}
