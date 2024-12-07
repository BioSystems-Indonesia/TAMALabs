package tcp

import (
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/usecase/hl_seven"
)

// HlSevenHandler is a struct that contains the handler of the REST server.
type HlSevenHandler struct {
	HLSevenUsecase *hl_seven.Usecase
}

// NewHlSevenHandler creates a new instance of HlSevenHandler.
func NewHlSevenHandler(hLSevenUsecase *hl_seven.Usecase) *HlSevenHandler {
	return &HlSevenHandler{
		HLSevenUsecase: hLSevenUsecase,
	}
}

// ProcessORM is a function to process ORM message.
func (h *HlSevenHandler) ProcessORM(message string) (string, error) {
	parse, err := entity.ORMMessage(message).Parse()
	if err != nil {
		return "", err
	}
	return h.HLSevenUsecase.ProcessORM(parse)
}
