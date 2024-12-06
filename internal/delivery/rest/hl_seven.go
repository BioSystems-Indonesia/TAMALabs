package rest

import (
	"github.com/labstack/echo/v4"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/usecase/hl_seven"
)

// HlSevenHandler is a struct that contains the handler of the REST server.
type HlSevenHandler struct {
	HLSevenUsecase *hl_seven.Usecase
}

func NewHlSevenHandler(hLSevenUsecase *hl_seven.Usecase) *HlSevenHandler {
	return &HlSevenHandler{
		HLSevenUsecase: hLSevenUsecase,
	}
}

// SendORM is a function to send ORM message.
func (h *HlSevenHandler) SendORM(c echo.Context) error {
	var message entity.SendORMRequest
	if err := c.Bind(&message); err != nil {
		return c.JSON(400, "Failed to bind")
	}
	resp, err := h.HLSevenUsecase.SendORM(message)
	if err != nil {
		return err
	}
	return c.JSON(200, resp)
}
