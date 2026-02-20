package repositories

import (
	"context"
	"time"

	"github.com/BioSystems-Indonesia/integration-services-lis/internal/models"
	"gorm.io/gorm"
)

type TestTypeRepository interface {
	FindAll(ctx context.Context) ([]models.TestType, error)
	UpdateLastSync(ctx context.Context, testTypeIds []int, syncTime time.Time) error
}

type TestTypeRepositoryImpl struct {
	db *gorm.DB
}

func NewTestTypeRepository(db *gorm.DB) TestTypeRepository {
	return &TestTypeRepositoryImpl{db: db}
}

func (t TestTypeRepositoryImpl) FindAll(ctx context.Context) ([]models.TestType, error) {
	var testTypes []models.TestType

	err := t.db.WithContext(ctx).
		Where(`
			(
				last_sync IS NULL
				OR updated_at IS NULL
				OR updated_at > last_sync
			)
		`).
		Order("code ASC").
		Find(&testTypes).
		Error

	if err != nil {
		return nil, err
	}

	return testTypes, nil
}

func (t TestTypeRepositoryImpl) UpdateLastSync(ctx context.Context, testTypeIds []int, syncTime time.Time) error {
	if len(testTypeIds) == 0 {
		return nil
	}

	return t.db.WithContext(ctx).
		Model(&models.TestType{}).
		Where("id IN ?", testTypeIds).
		UpdateColumn("last_sync", syncTime).
		Error
}
