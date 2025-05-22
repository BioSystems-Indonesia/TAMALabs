package adminrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql"
	"gorm.io/gorm"
)

type AdminRepository struct {
	db  *gorm.DB
	cfg *config.Schema
}

func NewAdminRepository(db *gorm.DB, cfg *config.Schema) *AdminRepository {
	return &AdminRepository{db: db, cfg: cfg}
}

func (r AdminRepository) FindAll(
	ctx context.Context,
	req *entity.GetManyRequestAdmin,
) (entity.PaginationResponse[entity.Admin], error) {
	db := r.db.WithContext(ctx).Preload("Roles")
	db = sql.ProcessGetMany(db, req.GetManyRequest,
		sql.Modify{
			ProcessSearch: func(db *gorm.DB, query string) *gorm.DB {
				return db.Where("fullname like ? or email like ? or username like ?", query+"%", query+"%", query+"%")
			},
		})

	return sql.GetWithPaginationResponse[entity.Admin](db, req.GetManyRequest)
}

func (r AdminRepository) FindOne(id int64) (entity.Admin, error) {
	var admin entity.Admin
	err := r.db.Where("id = ?", id).Preload("Roles").First(&admin).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Admin{}, entity.ErrNotFound
	}

	if err != nil {
		return entity.Admin{}, fmt.Errorf("error finding admin: %w", err)
	}

	return admin, nil
}

func (r AdminRepository) Create(admin *entity.Admin) error {
	return r.db.Create(admin).Error
}

func (r AdminRepository) Update(admin *entity.Admin) error {
	res := r.db.Updates(admin).Error
	if res != nil {
		return fmt.Errorf("error updating admin: %w", res)
	}

	return nil
}

func (r AdminRepository) Delete(id int64) error {
	res := r.db.Delete(&entity.Admin{ID: id})
	if res.Error != nil {
		return fmt.Errorf("error deleting admin: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return entity.ErrNotFound
	}

	return nil
}

func (r AdminRepository) FindOneByEmail(ctx context.Context, email string) (entity.Admin, error) {
	var admin entity.Admin
	err := r.db.Where("email = ?", email).First(&admin).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Admin{}, entity.ErrNotFound
	}

	if err != nil {
		return entity.Admin{}, fmt.Errorf("error finding admin: %w", err)
	}

	return admin, nil
}

func (r AdminRepository) FindOneByUsername(ctx context.Context, username string) (entity.Admin, error) {
	var admin entity.Admin
	err := r.db.Where("username = ?", username).First(&admin).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Admin{}, entity.ErrNotFound
	}

	if err != nil {
		return entity.Admin{}, fmt.Errorf("error finding admin: %w", err)
	}

	return admin, nil
}
