package analyxpanca

import (
	"fmt"
	"strconv"
	"time"

	"github.com/kardianos/hl7/h231"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

func encodeToPID(in entity.Patient) *h231.PID {
	id := strconv.FormatInt(in.ID, 10)
	return &h231.PID{
		HL7:                   h231.HL7Name{},
		SetID:                 "1",
		PatientID:             h231.CX{ID: id},
		PatientIdentifierList: []h231.CX{{ID: id}},
		AlternatePatientID:    []h231.CX{},
		PatientName: []h231.XPN{{
			FamilyNameLastNamePrefix: in.LastName,
			GivenName:                in.FirstName,
		}},
		MotherSMaidenName: []h231.XPN{},
		DateTimeOfBirth:   in.Birthdate,
		Sex:               in.Sex.String(),
	}
}

func getPatientID(id h231.CX) int64 {
	patientID, err := strconv.Atoi(id.ID)
	if err != nil {
		return 0
	}
	return int64(patientID)
}

func getPatientName(names []h231.XPN) (string, string) {
	if len(names) == 0 {
		return "", ""
	}
	return names[0].GivenName, names[0].FamilyNameLastNamePrefix
}

func getPatientAddress(addresses []h231.XAD) string {
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

func mapPIDToPatientEntity(pid *h231.PID) entity.Patient {
	if pid == nil {
		return entity.Patient{}
	}

	firstName, lastName := getPatientName(pid.PatientName)

	return entity.Patient{
		ID:        getPatientID(pid.PatientID),
		FirstName: firstName,
		LastName:  lastName,
		Birthdate: pid.DateTimeOfBirth,
		Sex:       entity.PatientSex(pid.Sex),
		Location:  getPatientAddress(pid.PatientAddress),
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
}
