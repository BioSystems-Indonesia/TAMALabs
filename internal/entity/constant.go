package entity

const (
	// HeaderXTotalCount is a header for react admin. Expected value is a number with a total of the results
	HeaderXTotalCount = "X-Total-Count"
	// HeaderAuthorization is a header for auth. Expected value is a JWT token
	HeaderAuthorization = "Authorization"
)

const (
	// ContextKeyUser is a key for user in context
	ContextKeyUser = "user"
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

func (s SpecimenType) Name() string {
	switch s {
	case SpecimenTypeSER:
		return "Serum"
	case SpecimenTypeURI:
		return "Urine"
	case SpecimenTypePLM:
		return "Plasma"
	case SpecimenTypeWBL:
		return "Whole blood"
	case SpecimenTypeCSF:
		return "Cerebrospinal fluid"
	case SpecimenTypeLIQ:
		return "Biological liquid"
	case SpecimenTypeSEM:
		return "Semen"
	default:
		return ""
	}
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
