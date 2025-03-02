package unit

import (
	"context"
	"gorm.io/gorm"

	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql"
)

type Repository struct {
	DB  *gorm.DB
	cfg *config.Schema
}

// NewRepository creates a new test type repository
func NewRepository(db *gorm.DB, cfg *config.Schema) *Repository {
	return &Repository{DB: db, cfg: cfg}
}

func (r *Repository) FindAll(
	ctx context.Context, req *entity.UnitGetManyRequest,
) (entity.PaginationResponse[entity.Unit], error) {
	db := r.DB.WithContext(ctx)
	db = sql.ProcessGetMany(db, req.GetManyRequest,
		sql.Modify{
			ProcessSearch: func(db *gorm.DB, query string) *gorm.DB {
				return db.Where("key like ?", "%"+query+"%")
			},
		})

	return sql.GetWithPaginationResponse[entity.Unit](db, req.GetManyRequest)
}
