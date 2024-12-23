package tcp

import (
	"context"
	"log"

	"github.com/kardianos/hl7"
	"github.com/kardianos/hl7/h251"
	"github.com/oibacidem/lims-hl-seven/internal/usecase/analyzer"
)

// HlSevenHandler is a struct that contains the handler of the REST server.
type HlSevenHandler struct {
	AnalyzerUsecase *analyzer.Usecase
}

// NewHlSevenHandler creates a new instance of HlSevenHandler.
func NewHlSevenHandler(analyzerUsecase *analyzer.Usecase) *HlSevenHandler {
	return &HlSevenHandler{
		AnalyzerUsecase: analyzerUsecase,
	}
}

// HL7Handler handles the HL7 message.
func (h *HlSevenHandler) HL7Handler(ctx context.Context, message string) (string, error) {
	msgByte := []byte(message)
	d := hl7.NewDecoder(h251.Registry, nil)
	msg, err := d.Decode(msgByte)
	if err != nil {
		return "", err
	}

	switch m := msg.(type) {
	case h251.OUL_R22:
		data, err := MapOULR22ToEntity(&m)
		if err != nil {
			return "", err
		}
		err = h.AnalyzerUsecase.ProcessOULR22(ctx, data)
		if err != nil {
			return "", err
		}
	case h251.OUL_R21:
		log.Println(m)
	}

	if err != nil {
		return "", err
	}

	return "OBX processed", nil
}
