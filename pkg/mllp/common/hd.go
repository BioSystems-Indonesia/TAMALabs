package common

import (
	"github.com/kardianos/hl7/h231"
	"github.com/kardianos/hl7/h251"
)

func SimpleHD(id string) *h251.HD {
	return &h251.HD{
		HL7:             h251.HL7Name{},
		NamespaceID:     id,
		UniversalID:     "",
		UniversalIDType: "",
	}
}

func SimpleHD231(id string) *h231.HD {
	return &h231.HD{
		HL7:             h231.HL7Name{},
		NamespaceID:     id,
		UniversalID:     "",
		UniversalIDType: "",
	}
}
