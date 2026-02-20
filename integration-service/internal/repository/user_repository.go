package repositories

import (
	"context"
	"time"

	"github.com/BioSystems-Indonesia/integration-services-lis/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindAll(ctx context.Context) ([]models.Admin, error)
	FindById(ctx context.Context, userId int) (*models.Admin, error)
	UpdateLastSync(ctx context.Context, userIds []int, syncTime time.Time) error
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) FindAll(ctx context.Context) ([]models.Admin, error) {
	var admins []models.Admin

	err := r.db.WithContext(ctx).
		Preload("Roles").
		Where("username != ?", "admin").
		Where("is_active = ?", true).
		Where(`
			(
				last_sync IS NULL
				OR updated_at IS NULL
				OR updated_at > last_sync
			)
		`).
		Order("fullname ASC").
		Find(&admins).Error

	if err != nil {
		return nil, err
	}

	return admins, nil
}

func (r *UserRepositoryImpl) UpdateLastSync(ctx context.Context, userIds []int, syncTime time.Time) error {
	if len(userIds) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).
		Model(&models.Admin{}).
		Where("id IN ?", userIds).
		UpdateColumn("last_sync", syncTime).
		Error
}

func (r UserRepositoryImpl) FindById(ctx context.Context, userId int) (*models.Admin, error) {
	var admin models.Admin

	err := r.db.WithContext(ctx).
		Preload("Roles").
		Where("id = ?", userId).
		First(&admin).Error

	if err != nil {
		return nil, err
	}

	return &admin, nil
}
