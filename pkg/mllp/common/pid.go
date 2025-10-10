package common

import (
	"fmt"
	"strconv"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	"github.com/kardianos/hl7/h251"
)

func EncodeToPID(in entity.Patient) *h251.PID {
	id := strconv.FormatInt(in.ID, 10)
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

func getPatientID(id h251.CX) int64 {
	patientID, err := strconv.Atoi(id.IDNumber)
	if err != nil {
		return 0
	}
	return int64(patientID)
}

func getPatientName(names []h251.XPN) (string, string) {
	if len(names) == 0 {
		return "", ""
	}
	return names[0].GivenName, names[0].FamilyName
}

func getPatientAddress(addresses []h251.XAD) string {
	if len(addresses) == 0 {
		return ""
	}
	addr := addresses[0]
	return fmt.Sprintf("%s %s %s %s %s",
		addr.StreetAddress,
		addr.City,
		addr.StateOrProvince,
		addr.ZipOrPostalCode,
		addr.Country)
}

func MapPIDToPatientEntity(pid *h251.PID) entity.Patient {
	if pid == nil {
		return entity.Patient{}
	}

	firstName, lastName := getPatientName(pid.PatientName)

	return entity.Patient{
		ID:        getPatientID(pid.PatientID),
		FirstName: firstName,
		LastName:  lastName,
		Birthdate: pid.DateTimeOfBirth,
		Sex:       entity.PatientSex(pid.AdministrativeSex),
		Location:  getPatientAddress(pid.PatientAddress),
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
}
