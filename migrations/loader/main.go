package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

func main() {
	models := []any{
		&entity.ObservationRequest{},
		&entity.ObservationResult{},
		&entity.Patient{},
		&entity.Specimen{},
		&entity.WorkOrder{},
		&entity.WorkOrderDevice{},
		&entity.Device{},
		&entity.Unit{},
		&entity.TestType{},
		&entity.Config{},
		&entity.TestTemplate{},
		&entity.TestTemplateTestType{},
		&entity.Admin{},
		&entity.Role{},
	}

	stmts, err := gormschema.New("sqlite").Load(models...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}

	_, err = io.WriteString(os.Stdout, stmts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to write gorm schema: %v\n", err)
		os.Exit(1)
	}
}
