package entity

import (
	"errors"
	"fmt"
	"github.com/oibacidem/lims-hl-seven/internal/constant"
	"github.com/oibacidem/lims-hl-seven/pkg/hl7"
	"reflect"
	"strings"
)

// MSA (Message Acknowledgment) segment (used in ACK responses)
type MSA struct {
	AcknowledgmentCode string `json:"acknowledgment_code" hl7:"1"` // AA = Application Accept, AE = Application Error, AR = Application Reject
	MessageControlID   string `json:"message_control_id" hl7:"2"`
	TextMessage        string `json:"text_message" hl7:"3"`
}

type ACKMessage string

// ACK represents the HL7 acknowledgment (response)
type ACK struct {
	MSH MSH `json:"msh"`
	MSA MSA `json:"msa"`
}

// Parse parses an HL7 ACK response message into an ACK struct.
func (a ACKMessage) Parse() (ACK, error) {
	result := ACK{}
	lines := strings.Split(strings.TrimSpace(string(a)), "\r")
	if len(lines) == 0 {
		return ACK{}, errors.New("invalid HL7 message: no segments found")
	}
	for _, line := range lines {
		fields := strings.Split(line, "|")
		segmentType := fields[0]

		// Identify the struct to populate based on the segment type
		var segmentStruct interface{}
		switch segmentType {
		case constant.MSH:
			segmentStruct = reflect.ValueOf(result).Elem().FieldByName(constant.MSH).Addr().Interface()
		case constant.MSA:
			segmentStruct = &result.MSA

		default:
			continue // Ignore unrecognized segments
		}

		if segmentStruct != nil {
			if err := hl7.PopulateStruct(segmentStruct, fields); err != nil {
				return ACK{}, fmt.Errorf("failed to parse segment %s: %w", segmentType, err)
			}
		}
	}

	return result, nil
}
