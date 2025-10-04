package cbs400

import (
	"github.com/oibacidem/lims-hl-seven/internal/usecase"
	"go.bug.st/serial"
)

type Handler struct {
	analyzerUseCase usecase.Analyzer
}

func NewHandler(analyzerUsecase usecase.Analyzer) *Handler {
	return &Handler{analyzerUseCase: analyzerUsecase}
}

func (h *Handler) Handle(port serial.Port) {

}
