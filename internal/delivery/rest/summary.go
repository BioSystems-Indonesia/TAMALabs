package rest

import (
	"net/http"

	summary_uc "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/summary"
	"github.com/labstack/echo/v4"
)

type SummaryHandler struct {
	summaryUc *summary_uc.SummaryUseCase
}

func NewSummaryHandler(summaryUc *summary_uc.SummaryUseCase) *SummaryHandler {
	return &SummaryHandler{summaryUc: summaryUc}
}

func (h *SummaryHandler) GetAllSummary(c echo.Context) error {
	resp := h.summaryUc.Summary(c.Request().Context())
	return c.JSON(http.StatusOK, resp)
}
