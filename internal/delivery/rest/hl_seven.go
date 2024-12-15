package rest

import (
	"github.com/oibacidem/lims-hl-seven/internal/usecase/analyzer"
)

// HlSevenHandler is a struct that contains the handler of the REST server.
type HlSevenHandler struct {
	HLSevenUsecase *analyzer.Usecase
}

func NewHlSevenHandler(hLSevenUsecase *analyzer.Usecase) *HlSevenHandler {
	return &HlSevenHandler{
		HLSevenUsecase: hLSevenUsecase,
	}
}
