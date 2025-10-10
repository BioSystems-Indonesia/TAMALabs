package main

import (
	"fmt"
	"io"
	"os"

	"log/slog"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
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
		&entity.SequenceDaily{},
	}

	stmts, err := gormschema.New("sqlite").Load(models...)
	if err != nil {
		slog.Error("failed to load gorm schema", "error", err)
		panic(fmt.Sprintf("failed to load gorm schema: %v", err))
	}

	_, err = io.WriteString(os.Stdout, stmts)
	if err != nil {
		slog.Error("failed to write gorm schema", "error", err)
		panic(fmt.Sprintf("failed to write gorm schema: %v", err))
	}
}
