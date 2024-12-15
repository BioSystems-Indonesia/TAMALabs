package rest

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

type FeatureListHandler struct {
	features map[string]entity.Tables
}

func NewFeatureListHandler() *FeatureListHandler {
	return &FeatureListHandler{features: entity.TableList}
}

func (f FeatureListHandler) registerFeatureListHandler(feature entity.Tables) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req entity.GetManyRequest
		if err := bindAndValidate(c, &req); err != nil {
			return handleError(c, err)
		}

		tables := slices.Clone(feature)
		if req.Query != "" {
			tables = tables.FilterName(req.Query)
		}

		if len(req.ID) > 0 {
			tables = tables.FilterID(req.ID)
		}

		c.Response().Header().Set(entity.HeaderXTotalCount, strconv.Itoa(len(tables)))

		if tables == nil {
			tables = []entity.Table{}
		}
		return c.JSON(http.StatusOK, tables)
	}
}

func (f FeatureListHandler) RegisterFeatureList(group *echo.Group) {
	for key, feature := range f.features {
		group.GET(fmt.Sprintf("/feature-list-%s", key), f.registerFeatureListHandler(feature))
	}
}
