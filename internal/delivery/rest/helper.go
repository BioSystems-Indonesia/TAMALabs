package rest

import (
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

func handleError(
	c echo.Context,
	err error,
	extraInfo ...map[string]interface{},
) error {
	code, payload := getErrorPayload(extraInfo, err, c)

	return c.JSON(code, payload)
}

func handleErrorSSE(
	c echo.Context,
	w http.ResponseWriter,
	err error,
	extraInfo ...map[string]interface{},
) error {
	code, payload := getErrorPayload(extraInfo, err, c)
	w.Write([]byte(fmt.Sprintf("event: error\ndata: %s\n\n", payload.Error)))

	return c.NoContent(code)
}

func getErrorPayload(extraInfo []map[string]interface{}, err error, c echo.Context) (int, entity.ErrorPayload) {
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

	var userErr *entity.UserError
	if errors.As(err, &userErr) {
		code = http.StatusBadRequest
	}

	payload := entity.ErrorPayload{
		Path:       fmt.Sprintf("%s %s", c.Request().Method, c.Path()),
		StatusCode: code,
		Error:      err.Error(),
		ExtraInfo:  extraInfoMap,
	}
	logError(payload)
	return code, payload
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

func successPaginationResponse[T any](c echo.Context, result entity.PaginationResponse[T]) error {
	c.Response().Header().Set(entity.HeaderXTotalCount, strconv.Itoa(int(result.Total)))
	return c.JSON(http.StatusOK, result.Data)
}

// createSSEWriter creates a new server send event writer.
func createSSEWriter(c echo.Context) (http.ResponseWriter, http.Flusher, error) {
	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set(echo.HeaderConnection, "keep-alive")
	c.Response().WriteHeader(http.StatusOK)

	sseWriter := c.Response().Writer
	flusher, ok := sseWriter.(http.Flusher)
	if !ok {
		return nil, nil, entity.ErrInternalServerError.WithInternal(errors.New("streaming unsupported"))
	}

	return sseWriter, flusher, nil
}
