package tcp

import (
	"bufio"
	"context"
	"io"
	"log"
	"net"
	"strings"

	"github.com/oibacidem/lims-hl-seven/internal/constant"
	"github.com/oibacidem/lims-hl-seven/pkg/mllp"
	"github.com/oibacidem/lims-hl-seven/pkg/server"
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

func RegisterRotes2(s server.TCPServer, handler *Handler) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Recovered in f", r)
			}
		}()

		for {
			c := s.GetClient()
			mc := mllp.NewClient(c)
			b, err := mc.ReadAll()
			if err != nil {
				if err != io.EOF {
					log.Println(err)
				}
			}
			res, err := handler.HL7Handler(context.Background(), string(b))
			if err != nil {
				log.Println(err)
			}
			mc.Write([]byte(res))

			c.Close()
		}
	}()
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
		resp, err := handler.HL7Handler(context.TODO(), msg)
		if err != nil {
			return "NAK: " + err.Error()
		}
		return resp
	}
	return "NAK: Invalid message"
}
