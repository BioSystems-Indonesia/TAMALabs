package ba400

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/kardianos/hl7"
	"github.com/kardianos/hl7/h251"
	"github.com/labstack/gommon/log"
	"github.com/oibacidem/lims-hl-seven/internal/constant"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

// SendToBA400 is a function to send message to BA400 for now its singleton view
// this is temporary function because it need device entity..
func SendToBA400(ctx context.Context, patients []entity.Patient, device entity.Device) error {
	encoder := hl7.NewEncoder(&hl7.EncodeOption{
		TrimTrailingSeparator: true,
	})

	buf := bytes.Buffer{}
	for _, p := range patients {
		o := NewOML_O33(p)

		b, err := encoder.Encode(o)
		if err != nil {
			return fmt.Errorf("failed to encode oml_33: %w", err)
		}

		buf.Write(b)
		buf.Write([]byte{constant.FileSeparator, constant.CarriageReturn})
		buf.Write([]byte{constant.VerticalTab})
	}

	sender := Sender{
		host:     fmt.Sprintf("%s:%d", device.IPAddress, device.Port),
		deadline: time.Second * 120,
	}
	messageToSend := buf.Bytes()
	resp, err := sender.SendRaw(messageToSend)
	if err != nil {
		log.Errorj(map[string]interface{}{
			"message": "sending to BA400 failed",
			"raw":     string(messageToSend),
			"resp":    string(resp),
		})
		return fmt.Errorf("failed to send raw: %w", err)
	}

	log.Infoj(map[string]interface{}{
		"message": "sending to BA400",
		"raw":     string(messageToSend),
	})

	err = receiveResponse(resp)
	if err != nil {
		return fmt.Errorf("failed to receive response: %w", err)
	}

	return nil
}

// SendToBA400 is a function to send message to BA400 for now its singleton view
// this is temporary function because it need device entity..
func CancelOrderBA400(ctx context.Context, patients []entity.Patient, device entity.Device) error {
	encoder := hl7.NewEncoder(&hl7.EncodeOption{
		TrimTrailingSeparator: true,
	})

	buf := bytes.Buffer{}
	for _, p := range patients {
		o := NewOML_O33(p)

		b, err := encoder.Encode(o)
		if err != nil {
			return fmt.Errorf("failed to encode oml_33: %w", err)
		}

		buf.Write(b)
		buf.Write([]byte{constant.FileSeparator, constant.CarriageReturn})
		buf.Write([]byte{constant.VerticalTab})
	}

	sender := Sender{
		host:     fmt.Sprintf("%s:%d", device.IPAddress, device.Port),
		deadline: time.Second * 120,
	}
	messageToSend := buf.Bytes()
	resp, err := sender.SendRaw(messageToSend)
	if err != nil {
		log.Errorj(map[string]interface{}{
			"message": "sending to BA400 failed",
			"raw":     string(messageToSend),
			"resp":    string(resp),
		})
		return fmt.Errorf("failed to send raw: %w", err)
	}

	log.Infoj(map[string]interface{}{
		"message": "sending to BA400",
		"raw":     string(messageToSend),
	})

	err = receiveResponse(resp)
	if err != nil {
		return fmt.Errorf("failed to receive response: %w", err)
	}

	return nil
}

// receiveResponse is a function to receive response from BA400
// Resp example
// MSH|^~\\\u0026|BA200|Biosystems|Host|Host provider|20241223163505||ORL^O34^ORL_O34|ec3b41a9-77e3-4fd3-a2ca-f8f760dbda47|P|2.5.1|||ER|NE||UNICODE UTF-8|||LAB-28^IHE\rMSA|AA|939b894f-a10a-4b35-9f82-95de095cc0c4\r
func receiveResponse(resp []byte) error {
	decoder := hl7.NewDecoder(h251.Registry, nil)

	msg, err := decoder.Decode(resp)
	if err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	switch m := msg.(type) {
	case h251.ORL_O34:
		log.Infoj(map[string]interface{}{
			"message": "receive response",
			"resp":    fmt.Sprintf("%v", m),
		})
		s := msg.(h251.ORL_O34)

		switch s.MSA.AcknowledgmentCode {
		case h251.ID(constant.ApplicationAccept), h251.ID(constant.CommitAccept):
			return handleAccept(s)
		default:
			return fmt.Errorf("Got failed or reject acknowledgment code: %s", s.MSA.AcknowledgmentCode)
		}
	default:
		return fmt.Errorf("unsupported message type: %T", m)
	}
}

func handleAccept(s h251.ORL_O34) error {
	return nil
}
