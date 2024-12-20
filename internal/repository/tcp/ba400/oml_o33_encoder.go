package ba400

import (
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/kardianos/hl7/h251"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

func NewOML_O33(patient entity.Patient, sepecimen []entity.Specimen, observationRequest []entity.ObservationRequest) h251.OML_O33 {
	msgControlID := uuid.New()
	date := time.Now()
	return h251.OML_O33{
		MSH:      NewOML_O33_MSH(msgControlID.String(), date),
		Patient:  NewOML_O33_Patient(patient),
		Specimen: NewOML_O33_Specimens(sepecimen, observationRequest, date),
	}
}

func NewOML_O33_Specimens(sepeciments []entity.Specimen, obr []entity.ObservationRequest, date time.Time) []h251.OML_O33_Specimen {
	var specimens []h251.OML_O33_Specimen
	for i, s := range sepeciments {
		specimens = append(specimens, NewOML_O33_Specimen(i+1, s, obr, date))
	}
	return specimens
}

func NewOML_O33_Specimen(index int, s entity.Specimen, obr []entity.ObservationRequest, date time.Time) h251.OML_O33_Specimen {
	var orders = []h251.OML_O33_Order{}
	for i, o := range obr {
		orders = append(orders, NewOrder(strconv.Itoa(i+1), date, o.TestCode))
	}

	return h251.OML_O33_Specimen{
		HL7: h251.HL7Name{},
		SPM: &h251.SPM{

			HL7:   h251.HL7Name{},
			SetID: strconv.Itoa(index),
			SpecimenID: &h251.EIP{
				PlacerAssignedIdentifier: &h251.EI{
					EntityIdentifier: s.Barcode,
				},
			},
			SpecimenParentIDs: []h251.EIP{},
			SpecimenType:      Serum,
			SpecimenRole:      []h251.CWE{{Identifier: "P"}},
		},
		Order: orders,
	}
}

func SimpleHD(id string) *h251.HD {
	return &h251.HD{
		HL7:             h251.HL7Name{},
		NamespaceID:     id,
		UniversalID:     "",
		UniversalIDType: "",
	}
}

func NewOML_O33_MSH(id string, date time.Time) *h251.MSH {
	return &h251.MSH{
		HL7:                                 h251.HL7Name{},
		FieldSeparator:                      "|",
		EncodingCharacters:                  "^~\\&",
		SendingApplication:                  SimpleHD("BioLIS"),
		SendingFacility:                     SimpleHD("Lab1"),
		ReceivingApplication:                SimpleHD("BA200"),
		ReceivingFacility:                   SimpleHD("Lab1"),
		DateTimeOfMessage:                   date,
		Security:                            "",
		MessageType:                         OML_O33_MessageType,
		MessageControlID:                    id,
		ProcessingID:                        h251.PT{ProcessingID: "P"},
		VersionID:                           h251.VID{VersionID: "2.5.1"},
		SequenceNumber:                      "",
		ContinuationPointer:                 "",
		AcceptAcknowledgmentType:            "ER",
		ApplicationAcknowledgmentType:       "AL",
		CountryCode:                         "ID",
		CharacterSet:                        []string{"UNICODE UTF-8"},
		PrincipalLanguageOfMessage:          &h251.CE{},
		AlternateCharacterSetHandlingScheme: "",
		MessageProfileIdentifier: []h251.EI{
			{
				HL7:              h251.HL7Name{},
				EntityIdentifier: "LAB-28",
				NamespaceID:      "IHE",
				UniversalID:      "",
				UniversalIDType:  "",
			},
		},
	}
}

func NewOML_O33_Patient(patient entity.Patient) *h251.OML_O33_Patient {
	EncodeToPID(patient)
	return &h251.OML_O33_Patient{
		HL7: h251.HL7Name{},
		PID: EncodeToPID(patient),
	}
}

func NewOBRUniversalID(id string) h251.CE {
	return h251.CE{
		Identifier:         id,
		Text:               id,
		NameOfCodingSystem: "BA200",
	}
}

func NewOrder(setID string, date time.Time, testID string) h251.OML_O33_Order {
	return h251.OML_O33_Order{
		ORC: &h251.ORC{
			OrderControl:          "NW",
			DateTimeOfTransaction: date,
		},
		ObservationRequest: &h251.OML_O33_ObservationRequest{
			OBR: &h251.OBR{
				SetID:                      setID,
				PlacerOrderNumber:          &h251.EI{EntityIdentifier: "1"},
				UniversalServiceIdentifier: NewOBRUniversalID(testID),
			},
		},
	}
}

var OML_O33_MessageType h251.MSG = h251.MSG{
	HL7:              h251.HL7Name{},
	MessageCode:      "OML",
	TriggerEvent:     "O33",
	MessageStructure: "OML_O33",
}

var Serum h251.CWE = h251.CWE{
	Identifier:         "SER",
	Text:               "Serum",
	NameOfCodingSystem: "HL70369",
}
