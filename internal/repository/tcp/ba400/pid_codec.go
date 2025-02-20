package ba400

import (
	"fmt"

	"github.com/kardianos/hl7/h251"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

func EncodeToPID(in entity.Patient) *h251.PID {
	id := fmt.Sprintf("%d-%s %s", in.ID, in.FirstName, in.LastName)
	return &h251.PID{
		HL7:                   h251.HL7Name{},
		SetID:                 "1",
		PatientID:             h251.CX{IDNumber: id},
		PatientIdentifierList: []h251.CX{{IDNumber: id}},
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
