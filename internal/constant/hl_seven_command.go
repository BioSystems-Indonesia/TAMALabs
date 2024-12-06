package constant

const (
	// MSH is the message header segment
	// The MSH segment is used to transmit message header information.
	// The MSH segment is required in all HL7 messages.
	MSH = "MSH"

	// MSA is the message acknowledgment segment
	// The MSA segment is used to transmit message acknowledgment information.
	// The MSA segment is required in all ACK messages.
	MSA = "MSA"

	// ACK is the acknowledgment message
	// The ACK message is used to acknowledge the receipt of an HL7 message.
	ACK = "ACK"

	// PID is the patient identification segment
	// The PID segment is used to transmit patient identification information.
	// The PID segment is required in all ORM messages.
	PID = "PID"

	// ORC is the common order segment
	// The ORC segment is used to transmit information about an order, including placing, canceling, and filling an order.
	// The ORC segment is required in all ORM messages.
	ORC = "ORC"

	// OBR is the observation request segment
	// The OBR segment is used to transmit information about a test or observation request.
	// The OBR segment is required in all ORM messages.
	OBR = "OBR"

	// ORM is the order message
	// The ORM message is used to transmit an order for a test or observation.
	ORM = "ORM"
)
