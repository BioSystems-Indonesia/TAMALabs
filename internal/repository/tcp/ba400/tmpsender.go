package ba400

import (
	"fmt"
	"log"

	"github.com/kardianos/hl7"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

// SendToBA400 is a function to send message to BA400 for now its singleton view
// this is temporary function because it need device entity..
func SendToBA400(patient entity.Patient, sepecimen []entity.Specimen, observationRequest []entity.ObservationRequest) error {
	o := NewOML_O33(patient, sepecimen, observationRequest)
	encoder := hl7.NewEncoder(&hl7.EncodeOption{
		TrimTrailingSeparator: true,
	})
	b, err := encoder.Encode(o)
	if err != nil {
		return fmt.Errorf("failed to encode oml_33: %w", err)
	}
	log.Print(string(b))

	sender := Sender{
		host: "192.168.1.11:2050",
	}

	_, err = sender.SendRaw(b)
	if err != nil {
		return fmt.Errorf("failed to send raw: %w", err)
	}

	// TODO check response
	return nil
}
