package constant

const (
	// MSH (message header) is the first segment of an HL7 message
	MSH = "MSH"

	// ThisApplication are name of this application that will be used in MSH-3
	ThisApplication = "LIS"

	// ThisFacility are name of this facility that will be used in MSH-4
	ThisFacility = "Lab01"
)

type AcknowledgmentCode string

const (
	ApplicationAccept AcknowledgmentCode = "AA"
	ApplicationError  AcknowledgmentCode = "AE"
	ApplicationReject AcknowledgmentCode = "AR"
	CommitAccept      AcknowledgmentCode = "CA"
	CommitError       AcknowledgmentCode = "CE"
)

// Define a custom type for ResultStatus.
type ResultStatus string

// Define constants for OBX Record Status using the ResultStatus type.
const (
	// Record coming over is a correction and thus replaces a final result.
	ResultStatusCorrection ResultStatus = "C"
	// Deletes the OBX record.
	ResultStatusDelete ResultStatus = "D"
	// Final results; Can only be changed with a corrected result.
	ResultStatusFinal ResultStatus = "F"
	// Specimen in lab; results pending.
	ResultStatusSpecimenPending ResultStatus = "I"
	// Not asked; used to affirmatively document that the observation was not sought.
	ResultStatusNotAsked ResultStatus = "N"
	// Order detail description only (no result).
	ResultStatusOrderDetail ResultStatus = "O"
	// Preliminary results.
	ResultStatusPreliminary ResultStatus = "P"
	// Results entered -- not verified.
	ResultStatusNotVerified ResultStatus = "R"
	// Partial results.
	ResultStatusPartial ResultStatus = "S"
	// Results status change to final without retransmitting preliminary results.
	ResultStatusFinalStatus ResultStatus = "U"
	// Post original as wrong, e.g., transmitted for wrong patient.
	ResultStatusWrongPatient ResultStatus = "W"
	// Results cannot be obtained for this observation.
	ResultStatusCannotObtain ResultStatus = "X"
)

type OrderControlNode string

const (
	// Order/service refill request approval.
	OrderControlNodeAF OrderControlNode = "AF"
	// Cancel order/service request.
	OrderControlNodeCA OrderControlNode = "CA"
	// Child order/service.
	OrderControlNodeCH OrderControlNode = "CH"
	// Combined result.
	OrderControlNodeCN OrderControlNode = "CN"
	// Canceled as requested.
	OrderControlNodeCR OrderControlNode = "CR"
	// Discontinue order/service request.
	OrderControlNodeDC OrderControlNode = "DC"
	// Data errors.
	OrderControlNodeDE OrderControlNode = "DE"
	// Order/service refill request denied.
	OrderControlNodeDF OrderControlNode = "DF"
	// Discontinued as requested.
	OrderControlNodeDR OrderControlNode = "DR"
	// Order/service refilled, unsolicited.
	OrderControlNodeFU OrderControlNode = "FU"
	// Hold order request.
	OrderControlNodeHD OrderControlNode = "HD"
	// On hold as requested.
	OrderControlNodeHR OrderControlNode = "HR"
	// Link order/service to patient care problem or goal.
	OrderControlNodeLI OrderControlNode = "LI"
	// Number assigned.
	OrderControlNodeNA OrderControlNode = "NA"
	// New order/service.
	OrderControlNodeNW OrderControlNode = "NW"
	// Order/service canceled.
	OrderControlNodeOC OrderControlNode = "OC"
	// Order/service discontinued.
	OrderControlNodeOD OrderControlNode = "OD"
	// Order/service released.
	OrderControlNodeOE OrderControlNode = "OE"
	// Order/service refilled as requested.
	OrderControlNodeOF OrderControlNode = "OF"
	// Order/service held.
	OrderControlNodeOH OrderControlNode = "OH"
	// Order/service accepted & OK.
	OrderControlNodeOK OrderControlNode = "OK"
	// Notification of order for outside dispense.
	OrderControlNodeOP OrderControlNode = "OP"
	// Released as requested.
	OrderControlNodeOR OrderControlNode = "OR"
	// Parent order/service.
	OrderControlNodePA OrderControlNode = "PA"
	// Previous Results with new order/service.
	OrderControlNodePR OrderControlNode = "PR"
	// Notification of replacement order for outside dispense.
	OrderControlNodePY OrderControlNode = "PY"
	// Observations/Performed Service to follow.
	OrderControlNodeRE OrderControlNode = "RE"
	// Refill order/service request.
	OrderControlNodeRF OrderControlNode = "RF"
	// Release previous hold.
	OrderControlNodeRL OrderControlNode = "RL"
	// Replacement order.
	OrderControlNodeRO OrderControlNode = "RO"
	// Order/service replace request.
	OrderControlNodeRP OrderControlNode = "RP"
	// Replaced as requested.
	OrderControlNodeRQ OrderControlNode = "RQ"
	// Request received.
	OrderControlNodeRR OrderControlNode = "RR"
	// Replaced unsolicited.
	OrderControlNodeRU OrderControlNode = "RU"
	// Status changed.
	OrderControlNodeSC OrderControlNode = "SC"
	// Send order/service number.
	OrderControlNodeSN OrderControlNode = "SN"
	// Response to send order/service status request.
	OrderControlNodeSR OrderControlNode = "SR"
	// Send order/service status request.
	OrderControlNodeSS OrderControlNode = "SS"
	// Unable to accept order/service.
	OrderControlNodeUA OrderControlNode = "UA"
	// Unable to cancel.
	OrderControlNodeUC OrderControlNode = "UC"
	// Unable to discontinue.
	OrderControlNodeUD OrderControlNode = "UD"
	// Unable to refill.
	OrderControlNodeUF OrderControlNode = "UF"
	// Unable to put on hold.
	OrderControlNodeUH OrderControlNode = "UH"
	// Unable to replace.
	OrderControlNodeUM OrderControlNode = "UM"
	// Unlink order/service from patient care problem or goal.
	OrderControlNodeUN OrderControlNode = "UN"
	// Unable to release.
	OrderControlNodeUR OrderControlNode = "UR"
	// Unable to change.
	OrderControlNodeUX OrderControlNode = "UX"
	// Change order/service request.
	OrderControlNodeXO OrderControlNode = "XO"
	// Changed as requested.
	OrderControlNodeXR OrderControlNode = "XR"
	// Order/service changed, unsol.
	OrderControlNodeXX OrderControlNode = "XX"
)
