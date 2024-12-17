package tcp

import (
	"fmt"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"strconv"
	"time"

	"github.com/kardianos/hl7/h251"
)

func MapOULR22ToEntity(msg *h251.OUL_R22) (entity.OUL_R22, error) {
	// HL7Message Mapping (MSH Segment)
	hl7Message := entity.HL7Message{
		MessageControlID:     msg.MSH.MessageControlID,
		SendingApplication:   msg.MSH.SendingApplication.NamespaceID,
		SendingFacility:      msg.MSH.SendingFacility.NamespaceID,
		ReceivingApplication: msg.MSH.ReceivingApplication.NamespaceID,
		ReceivingFacility:    msg.MSH.ReceivingFacility.NamespaceID,
		MessageType:          msg.MSH.MessageType.MessageStructure,
		MessageVersion:       msg.MSH.VersionID.VersionID,
		MessageDate:          msg.MSH.DateTimeOfMessage.Format(time.RFC3339),
	}

	// Patient Mapping (PID Segment)
	patientID, err := strconv.Atoi(msg.Patient.PID.PatientID.IDNumber)
	if err != nil {
		patientID = 0
	}
	patient := entity.Patient{
		ID:        int64(patientID),
		FirstName: msg.Patient.PID.PatientName[0].GivenName,
		LastName:  msg.Patient.PID.PatientName[0].FamilyName,
		Birthdate: msg.Patient.PID.DateTimeOfBirth,
		Sex:       entity.PatientSex(msg.Patient.PID.AdministrativeSex),
		Location:  fmt.Sprintf("%s %s %s %s %s", msg.Patient.PID.PatientAddress[0].StreetAddress, msg.Patient.PID.PatientAddress[0].City, msg.Patient.PID.PatientAddress[0].StateOrProvince, msg.Patient.PID.PatientAddress[0].ZipOrPostalCode, msg.Patient.PID.PatientAddress[0].Country),
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	// Specimen Mapping (SPM Segment)
	specimen := entity.Specimen{
		HL7ID:        msg.Specimen[0].SPM.SpecimenID.PlacerAssignedIdentifier.EntityIdentifier,
		Type:         msg.Specimen[0].SPM.SpecimenType.Identifier,
		ReceivedDate: msg.Specimen[0].SPM.SpecimenReceivedDateTime,
	}

	// Observation Request Mapping (OBR Segment)
	order := msg.Specimen[0].Order[0]
	observationRequest := entity.ObservationRequest{
		ID:              0,
		SpecimenID:      0,
		OrderID:         order.OBR.PlacerOrderNumber.EntityIdentifier,
		TestCode:        order.OBR.UniversalServiceIdentifier.Identifier,
		TestDescription: order.OBR.UniversalServiceIdentifier.Text,
		RequestedDate:   order.OBR.RequestedDateTime,
		ResultStatus:    order.OBR.ResultStatus,
	}

	// Observations Mapping (OBX Segments)
	var observations []entity.Observation
	for _, v := range order.Result {
		obx := v.OBX
		observations = append(observations, entity.Observation{
			Code:           obx.ObservationIdentifier.Identifier,
			Description:    obx.ObservationIdentifier.Text,
			Value:          obx.ValueType,
			Unit:           obx.Units.Identifier,
			ReferenceRange: obx.ReferencesRange,
			Date:           obx.DateTimeOfTheObservation,
			AbnormalFlag:   obx.AbnormalFlags[0],
			Comments:       obx.ObservationResultStatus,
		})
	}

	// Return the composite struct
	return entity.OUL_R22{
		HL7Message:         hl7Message,
		Patient:            patient,
		Specimen:           specimen,
		ObservationRequest: observationRequest,
		Observations:       observations,
	}, nil
}
