package entity

const (
	// HeaderXTotalCount is a header for react admin. Expected value is a number with a total of the results
	HeaderXTotalCount = "X-Total-Count"
)

type SpecimenType string

const (
	// SpecimenTypeSER is a specimen type for serum.
	SpecimenTypeSER SpecimenType = "SER"
	// SpecimenTypeURI is a specimen type for urine.
	SpecimenTypeURI SpecimenType = "URI"
	// SpecimenTypePLM is a specimen type for plasma.
	SpecimenTypePLM SpecimenType = "PLM"
	// SpecimenTypeWBL is a specimen type for whole blood.
	SpecimenTypeWBL SpecimenType = "WBL"
	// SpecimenTypeCSF is a specimen type for cerebrospinal fluid.
	SpecimenTypeCSF SpecimenType = "CSF"
	// SpecimenTypeLIQ is a specimen type for biological liquid.
	SpecimenTypeLIQ SpecimenType = "LIQ"
	// SpecimenTypeSEM is a specimen type for semen.
	SpecimenTypeSEM SpecimenType = "SEM"
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
