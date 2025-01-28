package sql

import (
	"fmt"

	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Modify struct {
	ProcessSearch func(db *gorm.DB, query string) *gorm.DB
	ProcessID     func(db *gorm.DB, id []string) *gorm.DB
	// Sometimes, the table name need to be explicitly set. For example, in case of join to avoid ambiguous column name
	TableName string
}

func NoSearch(db *gorm.DB, query string) *gorm.DB {
	return db
}

func DefaultProcessID(db *gorm.DB, id []string, modify Modify) *gorm.DB {
	if len(id) > 0 {
		tableName := modify.TableName
		if tableName != "" {
			tableName += "."
		}
		query := fmt.Sprintf("%sid in (?)", tableName)
		db = db.Where(query, id)
	}

	return db
}

func ProcessGetMany(db *gorm.DB, req entity.GetManyRequest, modify Modify) *gorm.DB {
	if len(req.ID) > 0 {
		if modify.ProcessID != nil {
			db = modify.ProcessID(db, req.ID)
		} else if modify.ProcessID == nil {
			db = DefaultProcessID(db, req.ID, modify)
		}
	}

	if req.Sort != "" {
		db = db.Order(clause.OrderByColumn{
			Column: clause.Column{
				Name:  req.Sort,
				Table: modify.TableName,
			},
			Desc: req.IsSortDesc(),
		})
	}

	if modify.ProcessSearch != nil {
		var search string
		if req.Query != "" {
			search = req.Query
		}

		if req.Search != "" {
			search = req.Search
		}

		if search != "" {
			db = modify.ProcessSearch(db, search)
		}
	}

	return db
}

func GetWithPaginationResponse[T any](db *gorm.DB, req entity.GetManyRequest) (entity.PaginationResponse[T], error) {
	var m []T
	db = db.Model(&m)

	var count int64
	err := db.Count(&count).Error
	if err != nil {
		return entity.PaginationResponse[T]{}, fmt.Errorf("error counting model: %w", err)
	}

	if req.Start > 0 {
		offset := req.Start
		db = db.Offset(offset)
	}

	if req.End > 0 {
		limit := req.End - req.Start
		db = db.Limit(limit)
	}

	err = db.Find(&m).Error
	if err != nil {
		return entity.PaginationResponse[T]{}, fmt.Errorf("error finding model: %w", err)
	}

	return entity.PaginationResponse[T]{
		Data:  m,
		Total: count,
	}, nil
}
