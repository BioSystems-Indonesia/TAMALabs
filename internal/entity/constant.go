package entity

const (
	// HeaderXTotalCount is a header for react admin. Expected value is a number with a total of the results
	HeaderXTotalCount = "X-Total-Count"
)

type SpecimenType string

const (
	SpecimenTypeCSF SpecimenType = "CSF"
	SpecimenTypeLIQ SpecimenType = "LIQ"
	SpecimenTypePLM SpecimenType = "PLM"
	SpecimenTypeSEM SpecimenType = "SEM"
	SpecimenTypeSER SpecimenType = "SER"
	SpecimenTypeURI SpecimenType = "URI"
	SpecimenTypeWBL SpecimenType = "WBL"
)

// TODO: Change this to real data
func (s SpecimenType) Code() string {
	switch s {
	case SpecimenTypeSER:
		return "1"
	case SpecimenTypeURI:
		return "2"
	case SpecimenTypeCSF:
		return "3"
	case SpecimenTypeLIQ:
		return "4"
	case SpecimenTypePLM:
		return "5"
	case SpecimenTypeSEM:
		return "6"
	}
	return "9"
}

type Priority string

const (
	// As soon as possible (a priority lower than stat)
	PriorityA Priority = "A"
	// Preoperative (to be done prior to surgery)
	PriorityP Priority = "P"
	// Routine
	PriorityR Priority = "R"
	// Stat (do immediately)
	PriorityS Priority = "S"
	// Timing critical (do as near as possible to requested time)
	PriorityT Priority = "T"
)
