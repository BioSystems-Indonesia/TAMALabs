package rest

import (
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

func handleError(
	c echo.Context,
	err error,
	extraInfo ...map[string]interface{},
) error {
	var extraInfoMap map[string]interface{}
	if len(extraInfo) > 0 {
		extraInfoMap = extraInfo[0]
	}

	var httpErr *entity.HTTPError
	code := http.StatusInternalServerError
	if errors.As(err, &httpErr) {
		code = httpErr.Code
	}

	var echoErr *echo.HTTPError
	if errors.As(err, &echoErr) {
		code = echoErr.Code
	}

	payload := entity.ErrorPayload{
		Path:       fmt.Sprintf("%s %s", c.Request().Method, c.Path()),
		StatusCode: code,
		Error:      err.Error(),
		ExtraInfo:  extraInfoMap,
	}
	logError(payload)

	return c.JSON(code, payload)
}

func logError(
	payload entity.ErrorPayload,
) {
	extraInfo := payload.ExtraInfo
	logPayload := map[string]any{
		"path":        payload.Path,
		"error":       payload.Error,
		"status_code": payload.StatusCode,
	}
	maps.Copy(logPayload, extraInfo)

	slog.Error("http error", "data", logPayload)
}

func bindAndValidate(c echo.Context, v interface{}) error {
	if err := c.Bind(v); err != nil {
		return err
	}
	return c.Validate(v)
}
