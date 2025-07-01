package tcp

import (
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
	if patient == nil {
		return entity.Patient{}
	}

	return mapPIDToPatientEntity(patient.PID)
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
		specimenResult.Barcode = specimen.SPM.SpecimenID.PlacerAssignedIdentifier.EntityIdentifier
		specimenResult.Type = specimen.SPM.SpecimenType.Identifier
		specimenResult.ReceivedDate = specimen.SPM.SpecimenReceivedDateTime
	}

	return specimenResult
}

func mapOBRToObservationRequestEntity(obr *h251.OBR) entity.ObservationRequest {
	if obr == nil {
		return entity.ObservationRequest{}
	}

	TestCode := obr.UniversalServiceIdentifier.Identifier
	TestDescription := obr.UniversalServiceIdentifier.Text

	return entity.ObservationRequest{
		TestCode:        TestCode,
		TestDescription: TestDescription,
		RequestedDate:   obr.RequestedDateTime,
		ResultStatus:    obr.ResultStatus,
	}
}
