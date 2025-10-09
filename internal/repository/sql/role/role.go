package rolerepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/BioSystems-Indonesia/TAMALabs/config"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql"
	"gorm.io/gorm"
)

type RoleRepository struct {
	db  *gorm.DB
	cfg *config.Schema
}

func NewRoleRepository(db *gorm.DB, cfg *config.Schema) *RoleRepository {
	return &RoleRepository{db: db, cfg: cfg}
}

func (r RoleRepository) FindAll(
	ctx context.Context,
	req *entity.GetManyRequestRole,
) (entity.PaginationResponse[entity.Role], error) {
	db := r.db.WithContext(ctx)
	db = sql.ProcessGetMany(db, req.GetManyRequest,
		sql.Modify{
			ProcessSearch: func(db *gorm.DB, query string) *gorm.DB {
				return db.Where("fullname like ? or email like ?", query+"%", query+"%")
			},
		})

	return sql.GetWithPaginationResponse[entity.Role](db, req.GetManyRequest)
}

func (r RoleRepository) FindOne(id int64) (entity.Role, error) {
	var role entity.Role
	err := r.db.Where("id = ?", id).First(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Role{}, entity.ErrNotFound
	}

	if err != nil {
		return entity.Role{}, fmt.Errorf("error finding role: %w", err)
	}

	return role, nil
}
