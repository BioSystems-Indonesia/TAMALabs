package entity

import "fmt"

type LogStreamingResponse string

func NewLogStreamingResponse(message string) string {
	return fmt.Sprintf("data: %s\n\n", message)
}
