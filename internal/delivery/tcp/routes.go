package tcp

import (
	"bufio"
	"context"
	"io"
	"log"
	"strings"

	"github.com/oibacidem/lims-hl-seven/internal/constant"
	"github.com/oibacidem/lims-hl-seven/pkg/mllp"
	"github.com/oibacidem/lims-hl-seven/pkg/server"
)

type Handler struct {
	*HlSevenHandler
}

const EOT = "\r"

func Loop(s server.TCPServer, handler *Handler) {
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
