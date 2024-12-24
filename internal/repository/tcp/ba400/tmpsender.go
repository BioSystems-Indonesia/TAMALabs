package ba400

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/kardianos/hl7"
	"github.com/labstack/gommon/log"
	"github.com/oibacidem/lims-hl-seven/internal/constant"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

// SendToBA400 is a function to send message to BA400 for now its singleton view
// this is temporary function because it need device entity..
func SendToBA400(ctx context.Context, patients []entity.Patient) error {
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
		host:     "192.168.33.68:2050",
		deadline: time.Second * 5,
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
		"resp":    string(resp),
	})

	// TODO check response
	return nil
}
