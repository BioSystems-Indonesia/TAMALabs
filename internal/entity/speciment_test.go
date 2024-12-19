package entity

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestValidSpeciment(t *testing.T) {
	speciment := Specimen{
		Type:      "type",
		PatientID: 1,
	}

	validate := validator.New()
	if err := validate.Struct(speciment); err != nil {
		t.Error(err)
	}
}
