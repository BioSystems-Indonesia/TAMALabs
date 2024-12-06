package hl7

import (
	"fmt"
	"reflect"
	"strings"
)

// Serialize converts a struct into an HL7 message string.
func Serialize(data interface{}) (string, error) {
	var message strings.Builder

	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem() // Dereference the pointer if `data` is a pointer
	}

	for i := 0; i < v.NumField(); i++ {
		segment := v.Field(i).Interface()

		segmentMessage, err := serializeSegment(segment)
		if err != nil {
			return "", err
		}

		message.WriteString(segmentMessage + "\r")
	}

	return strings.TrimSuffix(message.String(), "\r"), nil
}

// serializeSegment converts a single segment struct into an HL7 segment string.
func serializeSegment(segment interface{}) (string, error) {
	v := reflect.ValueOf(segment)
	t := reflect.TypeOf(segment)

	if v.Kind() != reflect.Struct {
		return "", fmt.Errorf("segment must be a struct, got %v", v.Kind())
	}

	var fields []string

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		hl7Tag := t.Field(i).Tag.Get("hl7")

		if hl7Tag == "" {
			continue
		}

		// Convert the hl7 tag to an index
		index, err := stringToIndex(hl7Tag)
		if err != nil {
			continue
		}

		// Ensure the fields slice is large enough
		for len(fields) <= index {
			fields = append(fields, "")
		}

		// Set the field value in the correct position
		fields[index] = field.String()
	}

	// Construct the segment string
	segmentType := t.Name()
	return segmentType + "|" + strings.Join(fields, "|"), nil
}

// stringToIndex converts an HL7 tag (e.g., "3") to a zero-based index.
func stringToIndex(tag string) (int, error) {
	index := 0
	_, err := fmt.Sscanf(tag, "%d", &index)
	return index - 1, err
}
