package ba400

import (
	"strconv"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	"github.com/kardianos/hl7/h251"
)

func EncodeToPID(in entity.Patient) *h251.PID {
	id := strconv.FormatInt(in.ID, 10)
	return &h251.PID{
		HL7:       h251.HL7Name{},
		SetID:     "1",
		PatientID: h251.CX{IDNumber: id},
		// PatientID:             h251.CX{IDNumber: id},
		PatientIdentifierList: []h251.CX{{IDNumber: in.FirstName}},
		AlternatePatientID:    []h251.CX{},
		PatientName: []h251.XPN{{
			FamilyName: in.LastName,
			GivenName:  in.FirstName,
		}},
		MothersMaidenName: []h251.XPN{},
		DateTimeOfBirth:   in.Birthdate,
		AdministrativeSex: string(in.Sex),
	}
}
