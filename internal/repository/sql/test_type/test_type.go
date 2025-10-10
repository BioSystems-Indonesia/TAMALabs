package test_type

import (
	"context"

	"github.com/BioSystems-Indonesia/TAMALabs/config"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql"
	"gorm.io/gorm"
)

type Repository struct {
	DB  *gorm.DB
	cfg *config.Schema
}

// NewRepository creates a new test type repository
func NewRepository(db *gorm.DB, cfg *config.Schema) *Repository {
	return &Repository{DB: db, cfg: cfg}
}

// FindAll returns all test types.
func (r *Repository) FindAll(
	ctx context.Context, req *entity.TestTypeGetManyRequest,
) (entity.PaginationResponse[entity.TestType], error) {
	db := r.DB.Preload("Device")
	db = sql.ProcessGetMany(db, req.GetManyRequest, sql.Modify{
		ProcessSearch: func(db *gorm.DB, query string) *gorm.DB {
			return db.Where("name like ?", "%"+query+"%").
				Or("code like ?", "%"+query+"%").
				Or("(alias_code like ? AND alias_code != '')", "%"+query+"%")
		},
	})

	if req.Code != "" {
		db = db.Where("code like ? OR (alias_code like ? AND alias_code != '')", "%"+req.Code+"%", "%"+req.Code+"%")
	}

	if len(req.Categories) != 0 {
		db = db.Where("category in (?)", req.Categories)
	}

	if len(req.SubCategories) != 0 {
		db = db.Where("sub_category in (?)", req.SubCategories)
	}

	// Filter by device ID if provided
	if req.DeviceID != nil {
		db = db.Where("device_id = ?", *req.DeviceID)
	}

	resp, err := sql.GetWithPaginationResponse[entity.TestType](db, req.GetManyRequest)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (r *Repository) FindAllFilter(ctx context.Context) (entity.TestTypeFilter, error) {
	var categories []string
	db := r.DB.Distinct("category").Model(entity.TestType{}).
		Where("category <> ''").
		Pluck("category", &categories)

	var subCategories []string
	db = db.Distinct("sub_category").Model(entity.TestType{}).
		Where("sub_category <> ''").
		Pluck("sub_category", &subCategories)

	return entity.TestTypeFilter{
		Categories:    categories,
		SubCategories: subCategories,
	}, nil
}

func (r *Repository) FindOneByID(ctx context.Context, id int) (entity.TestType, error) {
	var data entity.TestType
	if err := r.DB.Preload("Device").First(&data, id).Error; err != nil {
		return entity.TestType{}, err
	}
	return data, nil
}

func (r *Repository) FindOneByCode(ctx context.Context, code string) (entity.TestType, error) {
	var data entity.TestType
	if err := r.DB.Preload("Device").Where("code = ?", code).First(&data).Error; err != nil {
		return entity.TestType{}, err
	}
	return data, nil
}

func (r *Repository) FindOneByAliasCode(ctx context.Context, aliasCode string) (entity.TestType, error) {
	var data entity.TestType
	if aliasCode == "" {
		return entity.TestType{}, gorm.ErrRecordNotFound
	}

	if err := r.DB.Preload("Device").Where("alias_code = ? AND alias_code != ''", aliasCode).First(&data).Error; err != nil {
		return entity.TestType{}, err
	}
	return data, nil
}

// FindOneByCodeAndSpecimenType finds test type by code and specimen type combination
func (r *Repository) FindOneByCodeAndSpecimenType(ctx context.Context, code string, specimenType string) (entity.TestType, error) {
	var data entity.TestType
	if err := r.DB.Preload("Device").Where("code = ? AND type LIKE ?", code, "%"+specimenType+"%").First(&data).Error; err != nil {

		return entity.TestType{}, err
	}
	return data, nil
}

// FindByCodeWithSpecimenTypes finds all test types with the same code but different specimen types
func (r *Repository) FindByCodeWithSpecimenTypes(ctx context.Context, code string) ([]entity.TestType, error) {
	var data []entity.TestType
	if err := r.DB.Preload("Device").Where("code = ?", code).Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// FindByDeviceID finds all test types associated with a specific device
func (r *Repository) FindByDeviceID(ctx context.Context, deviceID int) ([]entity.TestType, error) {
	var data []entity.TestType
	if err := r.DB.Preload("Device").Where("device_id = ?", deviceID).Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// FindUnassignedTestTypes finds all test types that are not assigned to any device
func (r *Repository) FindUnassignedTestTypes(ctx context.Context) ([]entity.TestType, error) {
	var data []entity.TestType
	if err := r.DB.Where("device_id IS NULL").Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (r *Repository) Create(ctx context.Context, req *entity.TestType) (entity.TestType, error) {
	if err := r.DB.Create(req).Error; err != nil {
		return entity.TestType{}, err
	}
	return *req, nil
}

func (r *Repository) Update(ctx context.Context, req *entity.TestType) (entity.TestType, error) {
	if err := r.DB.Save(req).Error; err != nil {
		return entity.TestType{}, err
	}
	return *req, nil
}

func (r *Repository) Delete(ctx context.Context, req *entity.TestType) (entity.TestType, error) {
	if err := r.DB.Delete(req).Error; err != nil {
		return entity.TestType{}, err
	}
	return *req, nil
}
