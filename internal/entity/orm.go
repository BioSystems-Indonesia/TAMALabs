package entity

import (
	"errors"
	"fmt"
	"github.com/oibacidem/lims-hl-seven/internal/constant"
	"github.com/oibacidem/lims-hl-seven/pkg/hl7"
	"reflect"
	"strings"
)

// ORM represents the ORM HL7 message (request)
type ORM struct {
	MSH MSH `json:"msh"`
	PID PID `json:"pid"`
	ORC ORC `json:"orc"`
	OBR OBR `json:"obr"`
}

type ORMMessage string

// SendORMRequest represents the ORM request.
type SendORMRequest struct {
	ORM ORM `json:"orm"`
}

type SendORMResponse struct {
	ACK ACK `json:"ack"`
}

func (o ORM) Serialize() string {
	message, err := hl7.Serialize(o)
	if err != nil {
		return ""
	}
	return message
}

// Parse parses an HL7 ORM request message into an ORM struct.
func (o ORMMessage) Parse() (ORM, error) {
	result := ORM{}
	lines := strings.Split(strings.TrimSpace(string(o)), "\r")
	if len(lines) == 0 {
		return ORM{}, errors.New("invalid HL7 message: no segments found")
	}
	for _, line := range lines {
		fields := strings.Split(line, "|")
		segmentType := fields[0]

		// Identify the struct to populate based on the segment type
		var segmentStruct interface{}
		switch segmentType {
		case constant.MSH:
			segmentStruct = reflect.ValueOf(result).Elem().FieldByName(constant.MSH).Addr().Interface()
		case constant.ORM:
			segmentStruct = &result.PID
		case constant.ORC:
			segmentStruct = &result.ORC
		case constant.OBR:
			segmentStruct = &result.OBR
		default:
			continue // Ignore unrecognized segments
		}

		if segmentStruct != nil {
			if err := hl7.PopulateStruct(segmentStruct, fields); err != nil {
				return ORM{}, fmt.Errorf("failed to parse segment %s: %w", segmentType, err)
			}
		}
	}

	return result, nil
}
