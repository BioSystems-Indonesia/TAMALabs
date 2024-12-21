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
	var observations []entity.Observation

	if specimen.OBX != nil {
		var observationResult []entity.ObservationResult
		for i := range specimen.OBX {
			observationResult = append(observationResult, entity.ObservationResult{
				Code:           specimen.OBX[i].ObservationIdentifier.Identifier,
				Description:    specimen.OBX[i].ObservationIdentifier.Text,
				Value:          mapObservationValueToValues(specimen.OBX[i].ObservationValue),
				Type:           specimen.OBX[i].ValueType,
				Unit:           specimen.OBX[i].Units.Identifier,
				ReferenceRange: specimen.OBX[i].ReferencesRange,
				Date:           specimen.OBX[i].DateTimeOfTheObservation,
				AbnormalFlag:   specimen.OBX[i].AbnormalFlags,
				Comments:       specimen.OBX[i].ObservationResultStatus,
			})
		}
		observations = append(observations, entity.Observation{
			Result: observationResult,
		})

	} else {
		for i := range specimen.Order {
			order := specimen.Order[i]

			// request
			obr := order.OBR
			observationRequest := entity.ObservationRequest{
				ID:              0,
				SpecimenID:      0,
				OrderID:         obr.PlacerOrderNumber.EntityIdentifier,
				TestCode:        obr.UniversalServiceIdentifier.Identifier,
				TestDescription: obr.UniversalServiceIdentifier.Text,
				RequestedDate:   obr.RequestedDateTime,
				ResultStatus:    obr.ResultStatus,
			}

			// result
			var observationResult []entity.ObservationResult
			for j := range order.Result {
				obx := order.Result[j].OBX
				observationResult = append(observationResult, entity.ObservationResult{
					Code:           obx.ObservationIdentifier.Identifier,
					Description:    obx.ObservationIdentifier.Text,
					Value:          mapObservationValueToValues(obx.ObservationValue),
					Type:           obx.ValueType,
					Unit:           obx.Units.Identifier,
					ReferenceRange: obx.ReferencesRange,
					Date:           obx.DateTimeOfTheObservation,
					AbnormalFlag:   obx.AbnormalFlags,
					Comments:       obx.ObservationResultStatus,
				})
			}
			observations = append(observations, entity.Observation{
				Request: observationRequest,
				Result:  observationResult,
			})
		}
	}

	return entity.Specimen{
		HL7ID:        specimen.SPM.SpecimenID.PlacerAssignedIdentifier.EntityIdentifier,
		Type:         specimen.SPM.SpecimenType.Identifier,
		ReceivedDate: specimen.SPM.SpecimenReceivedDateTime,
		Observation:  observations,
	}
}

func mapObservationValueToValues(values []h251.VARIES) []string {
	var results []string
	for i := range values {
		results = append(results, fmt.Sprintf("%v", values[i]))
	}
	return results
}
