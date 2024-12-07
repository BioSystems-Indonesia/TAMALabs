package tcp

import (
	"bufio"
	"github.com/oibacidem/lims-hl-seven/internal/constant"
	"log"
	"net"
	"strings"
)

type Handler struct {
	*HlSevenHandler
}

const EOT = "\r"

func RegisterRoutes(conn *net.TCPConn, handler *Handler) {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	var messageBuilder strings.Builder

	for {
		// Read the next line
		line, err := reader.ReadString('\n')
		if err != nil {
			if messageBuilder.Len() > 0 {
				// Process the remaining message if there's partial data
				processMessage(messageBuilder.String(), handler, writer)
			}
			log.Printf("Connection closed or read error: %v", err)
			return
		}

		line = strings.TrimSpace(line)

		// Accumulate lines into a complete message
		if strings.HasPrefix(line, constant.MSH) {
			// Start a new message
			messageBuilder.Reset()
		}
		messageBuilder.WriteString(line + "\r")

		// Check if this is the end of the message
		if strings.HasSuffix(line, EOT) || len(line) == 0 {
			// Process the complete message
			processMessage(messageBuilder.String(), handler, writer)
			messageBuilder.Reset()
		}
	}
}

// processMessage routes the message to the appropriate handler
func processMessage(message string, handler *Handler, writer *bufio.Writer) {
	log.Printf("Received complete message:\n%s", message)

	// Route the message
	response := routes(message, handler)

	// Send the response
	_, err := writer.WriteString(response + "\r")
	if err != nil {
		log.Printf("Error sending response: %v", err)
	}
	writer.Flush()
}

// routes processes the complete message and delegates it to the appropriate handler
func routes(msg string, handler *Handler) string {
	if strings.HasPrefix(msg, constant.MSH) {
		// Determine the message type from the MSH segment
		if strings.Contains(msg, constant.ORM) {
			_, err := handler.ProcessORM(msg)
			if err != nil {
				return "NAK: " + err.Error()
			}
			return "ACK: "
		} else if strings.Contains(msg, constant.ADT) {
			return ADTHandler(msg)
		}
		return "NAK: Unsupported HL7 message type"
	}
	return "NAK: Invalid HL7 message"
}

// ADTHandler processes ADT requests
func ADTHandler(msg string) string {
	log.Println("Processing ADT message...")
	return "ACK: ADT Message Processed"
}
