package ba400

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/kardianos/hl7"
	"github.com/kardianos/hl7/h251"
	"github.com/oibacidem/lims-hl-seven/internal/constant"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

const (
	timeout                 = 15 * time.Second
	batchObservationRequest = 30
	batchSend               = 3
)

// SendToBA400 is a function to send message to BA400 for now its singleton view
func SendToBA400(
	ctx context.Context, req *entity.SendPayloadRequest,
) error {
	encoder := hl7.NewEncoder(&hl7.EncodeOption{
		TrimTrailingSeparator: true,
	})

	splitObservationRequest := splitPatientsBasedOnObservationRequest(req.Patients)
	err := writeProgress(req, 0, entity.WorkOrderStreamingResponseStatusInProgress)
	if err != nil {
		return fmt.Errorf("failed to write progress: %w", err)
	}

	buf := bytes.Buffer{}
	for i, p := range splitObservationRequest {
		slog.Info(fmt.Sprintf("sending patient with len specimen %d", len(p.Specimen)))
		slog.Info(fmt.Sprintf("sending patient with len request %d", len(p.Specimen[0].ObservationRequest)))

		o := NewOML_O33(p, req.Device, req.Urgent)

		b, err := encoder.Encode(o)
		if err != nil {
			return fmt.Errorf("failed to encode oml_33: %w", err)
		}

		buf.Write(b)
		buf.Write([]byte{constant.FileSeparator, constant.CarriageReturn})
		buf.Write([]byte{constant.VerticalTab})

		if ((i+1)%batchSend == 0) || i == len(splitObservationRequest)-1 {
			sender := Sender{
				host:     fmt.Sprintf("%s:%d", req.Device.IPAddress, req.Device.Port),
				deadline: timeout,
			}
			messageToSend := buf.Bytes()
			slog.Info("sending to BA400",
				"raw", string(messageToSend),
				"i", i,
			)

			resp, err := sender.SendRaw(messageToSend)
			if err != nil {
				slog.Error("sending to BA400 failed",
					"raw", string(messageToSend),
					"resp", string(resp),
					"i", i,
				)
				return fmt.Errorf("failed to send raw: %w", err)
			}
			slog.Info("finish sending, waiting for response",
				"i", i,
			)

			err = receiveResponse(resp)
			if err != nil {
				return fmt.Errorf("failed to receive response: %w", err)
			}

			buf.Reset()
		}

		time.Sleep(time.Second)

		percentageComplete := float64(i+1) / float64(len(splitObservationRequest)) * 100
		err = writeProgress(req, percentageComplete, entity.WorkOrderStreamingResponseStatusInProgress)
		if err != nil {
			return fmt.Errorf("failed to write progress: %w", err)
		}
	}

	return nil
}

func writeProgress(req *entity.SendPayloadRequest, percentage float64, status entity.WorkOrderStreamingResponseStatus) error {
	if req.Writer == nil || req.Flusher == nil {
		return nil
	}

	_, err := req.Writer.Write([]byte(entity.NewWorkOrderStreamingResponse(percentage, status)))
	if err != nil {
		return fmt.Errorf("failed to write response: %w", err)
	}

	req.Flusher.Flush()

	return nil
}

func splitPatientsBasedOnObservationRequest(patients []entity.Patient) []entity.Patient {
	var splitPatient []entity.Patient
	for _, p := range patients {
		for _, s := range p.Specimen {
			for i := 0; i < len(s.ObservationRequest); i += batchObservationRequest {
				requests := []entity.ObservationRequest{}
				end := i + batchObservationRequest
				if end > len(s.ObservationRequest) {
					end = len(s.ObservationRequest)
				}

				r := s.ObservationRequest[i:end]
				requests = append(requests, r...)

				newSpecimen := s
				newSpecimen.ObservationRequest = requests
				newPatient := p
				newPatient.Specimen = []entity.Specimen{newSpecimen}
				splitPatient = append(splitPatient, newPatient)
			}
		}
	}
	return splitPatient
}

// receiveResponse is a function to receive response from BA400
// Resp example
// MSH|^~\\\u0026|BA200|Biosystems|Host|Host provider|20241223163505||ORL^O34^ORL_O34|ec3b41a9-77e3-4fd3-a2ca-f8f760dbda47|P|2.5.1|||ER|NE||UNICODE UTF-8|||LAB-28^IHE\rMSA|AA|939b894f-a10a-4b35-9f82-95de095cc0c4\r
func receiveResponse(resp []byte) error {
	slog.Info("receive response",
		"resp", resp,
	)

	decoder := hl7.NewDecoder(h251.Registry, nil)

	msg, err := decoder.Decode(resp)
	if err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	switch m := msg.(type) {
	case h251.ORL_O34:
		s := msg.(h251.ORL_O34)

		switch s.MSA.AcknowledgmentCode {
		case h251.ID(constant.ApplicationAccept), h251.ID(constant.CommitAccept):
			return handleAccept(s)
		default:
			return fmt.Errorf("got failed or reject acknowledgment code: %s", s.MSA.AcknowledgmentCode)
		}
	default:
		return fmt.Errorf("unsupported message type: %T", m)
	}
}

func handleAccept(_ h251.ORL_O34) error {
	return nil
}
