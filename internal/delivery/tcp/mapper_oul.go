package tcp

import (
	"fmt"
	"strconv"
	"time"

	"github.com/kardianos/hl7/h251"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

func MapOULR22ToEntity(msg *h251.OUL_R22) (entity.OUL_R22, error) {
	// MSH Mapping (MSH Segment)
	msh := mapMSHToEntity(msg.MSH)

	// Patient Mapping (PID Segment)
	var patient entity.Patient
	if msg.Patient != nil {
		patient = MapOULR22PatientToPatientEntity(msg.Patient)
	}

	// Specimens Mapping (SPM Segment)
	var specimens []entity.Specimen
	for i := range msg.Specimen {
		specimens = append(specimens, mapOULSpecimenToSpecimenEntity(msg.Specimen[i]))
	}

	// Return the composite struct
	return entity.OUL_R22{
		Msh:       msh,
		Patient:   patient,
		Specimens: specimens,
	}, nil
}

func MapOULR22PatientToPatientEntity(patient *h251.OUL_R22_Patient) entity.Patient {
	patientID, err := strconv.Atoi(patient.PID.PatientID.IDNumber)
	if err != nil {
		patientID = 0
	}
	return entity.Patient{
		ID:        int64(patientID),
		FirstName: patient.PID.PatientName[0].GivenName,
		LastName:  patient.PID.PatientName[0].FamilyName,
		Birthdate: patient.PID.DateTimeOfBirth,
		Sex:       entity.PatientSex(patient.PID.AdministrativeSex),
		Location: fmt.Sprintf("%s %s %s %s %s",
			patient.PID.PatientAddress[0].StreetAddress,
			patient.PID.PatientAddress[0].City,
			patient.PID.PatientAddress[0].StateOrProvince, patient.PID.PatientAddress[0].ZipOrPostalCode,
			patient.PID.PatientAddress[0].Country),
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
}

func mapOULSpecimenToSpecimenEntity(specimen h251.OUL_R22_Specimen) entity.Specimen {
	var (
		observationRequests []entity.ObservationRequest
		observationResults  []entity.ObservationResult
	)

	if specimen.OBX != nil {
		for i := range specimen.OBX {
			observationResults = append(observationResults, mapOBXToObservationResultEntity(&specimen.OBX[i]))
		}
	} else {
		for i := range specimen.Order {
			order := specimen.Order[i]

			// request
			observationRequests = append(observationRequests, mapOBRToObservationRequestEntity(order.OBR))

			// result
			for j := range order.Result {
				observationResults = append(observationResults, mapOBXToObservationResultEntity(order.Result[j].OBX))
			}
		}
	}

	specimenResult := entity.Specimen{
		ObservationRequest: observationRequests,
		ObservationResult:  observationResults,
	}

	if specimen.SPM != nil {
		specimenResult.HL7ID = specimen.SPM.SpecimenID.PlacerAssignedIdentifier.EntityIdentifier
		specimenResult.Type = specimen.SPM.SpecimenType.Identifier
		specimenResult.ReceivedDate = specimen.SPM.SpecimenReceivedDateTime
	}

	return specimenResult
}

func mapObservationValueToValues(values []h251.VARIES) []string {
	var results []string
	for i := range values {
		results = append(results, fmt.Sprintf("%v", values[i]))
	}
	return results
}

func mapOBRToObservationRequestEntity(obr *h251.OBR) entity.ObservationRequest {
	return entity.ObservationRequest{
		TestCode:        obr.UniversalServiceIdentifier.Identifier,
		TestDescription: obr.UniversalServiceIdentifier.Text,
		RequestedDate:   obr.RequestedDateTime,
		ResultStatus:    obr.ResultStatus,
	}
}

func mapOBXToObservationResultEntity(obx *h251.OBX) entity.ObservationResult {
	return entity.ObservationResult{
		Code:           obx.ObservationIdentifier.Identifier,
		Description:    obx.ObservationIdentifier.Text,
		Values:         mapObservationValueToValues(obx.ObservationValue),
		Type:           obx.ValueType,
		Unit:           obx.Units.Identifier,
		ReferenceRange: obx.ReferencesRange,
		Date:           obx.DateTimeOfTheObservation,
		AbnormalFlag:   obx.AbnormalFlags,
		Comments:       obx.ObservationResultStatus,
	}
}
