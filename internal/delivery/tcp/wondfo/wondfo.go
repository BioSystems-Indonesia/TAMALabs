package wondfo

import (
	"net"

	"github.com/oibacidem/lims-hl-seven/internal/usecase"
)

type Handler struct {
	analyzerUsecase usecase.Analyzer
}

func NewHandler(analyzerUsecase usecase.Analyzer) *Handler {
	return &Handler{
		analyzerUsecase: analyzerUsecase,
	}
}

func (h *Handler) Handle(conn *net.TCPConn) {}
