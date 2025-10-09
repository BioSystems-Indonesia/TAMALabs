package rest

import (
	"fmt"
	"net/http"
	"slices"
	"sort"
	"strconv"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	"github.com/labstack/echo/v4"
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

		tables, total := f.filterTables(feature, req)

		c.Response().Header().Set(entity.HeaderXTotalCount, strconv.Itoa(total))
		if tables == nil {
			tables = []entity.Table{}
		}
		return c.JSON(http.StatusOK, tables)
	}
}

func (f FeatureListHandler) filterTables(feature entity.Tables, req entity.GetManyRequest) (entity.Tables, int) {
	tables := slices.Clone(feature)

	if req.Query != "" {
		tables = tables.FilterName(req.Query)
	}

	if req.Sort != "" {
		sort.Slice(tables, func(i, j int) bool {
			if req.Sort == "id" {
				if req.IsSortDesc() {
					return tables[i].ID > tables[j].ID
				}

				return tables[i].ID < tables[j].ID
			}

			if req.Sort == "name" {
				if req.IsSortDesc() {
					return tables[i].Name > tables[j].Name
				}

				return tables[i].Name < tables[j].Name
			}

			return tables[i].Name < tables[j].Name
		})
	}

	if len(req.ID) > 0 {
		tables = tables.FilterID(req.ID)
	}

	total := len(tables)
	if req.Start != 0 || req.End != 0 {
		if len(tables) < req.End {
			req.End = len(tables)
		}

		tables = tables[req.Start:req.End]
	}

	return tables, total
}

func (f FeatureListHandler) RegisterFeatureList(group *echo.Group) {
	for key, feature := range f.features {
		group.GET(fmt.Sprintf("/feature-list-%s", key), f.registerFeatureListHandler(feature))
	}
}
