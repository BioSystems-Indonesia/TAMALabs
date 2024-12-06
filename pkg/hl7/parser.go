package hl7

import (
	"reflect"
	"strconv"
)

// PopulateStruct populates a struct's fields using reflection and numeric hl7 tags.
func PopulateStruct(segmentStruct interface{}, fields []string) error {
	v := reflect.ValueOf(segmentStruct).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if !field.CanSet() {
			continue
		}

		// Get the hl7 tag
		hl7Tag := t.Field(i).Tag.Get("hl7")
		if hl7Tag == "" {
			continue
		}

		// Convert the hl7 tag to an integer index
		index, err := strconv.Atoi(hl7Tag)
		if err != nil || index < 1 || index > len(fields) {
			continue
		}

		// Set the field value
		field.SetString(fields[index-1])
	}

	return nil
}
