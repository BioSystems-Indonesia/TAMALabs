package test_template

import (
	"context"

	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql"
	"github.com/oibacidem/lims-hl-seven/internal/util"
	"gorm.io/gorm"
)

type Repository struct {
	DB  *gorm.DB
	cfg *config.Schema
}

// NewRepository creates a new test type repository.
func NewRepository(db *gorm.DB, cfg *config.Schema) *Repository {
	return &Repository{DB: db, cfg: cfg}
}

// FindAll returns all test types.
func (r *Repository) FindAll(
	_ context.Context,
	req *entity.TestTemplateGetManyRequest,
) (entity.PaginationResponse[entity.TestTemplate], error) {
	db := r.DB
	sql.ProcessGetMany(db, req.GetManyRequest, sql.Modify{
		ProcessSearch: func(db *gorm.DB, query string) *gorm.DB {
			return db.Where("name like ?", "%"+query+"%").
				Or("description like ?", "%"+query+"%")
		}})

	resp, err := sql.GetWithPaginationResponse[entity.TestTemplate](db, req.GetManyRequest)
	if err != nil {
		return entity.PaginationResponse[entity.TestTemplate]{}, err
	}

	for i := range resp.Data {
		resp.Data[i].TestTypeID = util.Map(resp.Data[i].TestType, func(t entity.TestType) int {
			return t.ID
		})
	}

	return resp, nil
}

func (r *Repository) FindOneByID(ctx context.Context, id int) (entity.TestTemplate, error) {
	var data entity.TestTemplate
	if err := r.DB.Preload("TestType").First(&data, id).Error; err != nil {
		return entity.TestTemplate{}, err
	}

	data.TestTypeID = util.Map(data.TestType, func(t entity.TestType) int {
		return t.ID
	})

	return data, nil
}

func (r *Repository) Create(ctx context.Context, req *entity.TestTemplate) (entity.TestTemplate, error) {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(req).Error; err != nil {
			return err
		}

		err := r.updateRelation(tx, req)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return entity.TestTemplate{}, err
	}

	return *req, nil
}

func (r *Repository) Update(ctx context.Context, req *entity.TestTemplate) (entity.TestTemplate, error) {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(req).Error; err != nil {
			return err
		}

		err := r.updateRelation(tx, req)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return entity.TestTemplate{}, err
	}

	return *req, nil
}

func (r *Repository) updateRelation(tx *gorm.DB, req *entity.TestTemplate) error {
	err := tx.Delete(&entity.TestTemplateTestType{}, "test_template_id = ?", req.ID).Error
	if err != nil {
		return err
	}

	for _, testTypeID := range req.TestTypeID {
		if err := tx.Create(&entity.TestTemplateTestType{
			TestTemplateID: req.ID,
			TestTypeID:     testTypeID,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, req *entity.TestTemplate) (entity.TestTemplate, error) {
	if err := r.DB.Delete(req).Error; err != nil {
		return entity.TestTemplate{}, err
	}
	return *req, nil
}
