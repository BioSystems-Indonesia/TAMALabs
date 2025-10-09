package devicerepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/BioSystems-Indonesia/TAMALabs/config"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql"
	"gorm.io/gorm"
)

type DeviceRepository struct {
	db  *gorm.DB
	cfg *config.Schema
}

func NewDeviceRepository(db *gorm.DB, cfg *config.Schema) *DeviceRepository {
	return &DeviceRepository{db: db, cfg: cfg}
}

func (r DeviceRepository) FindAll(
	ctx context.Context,
	req *entity.GetManyRequestDevice,
) (entity.PaginationResponse[entity.Device], error) {
	db := r.db.WithContext(ctx)
	db = sql.ProcessGetMany(db, req.GetManyRequest,
		sql.Modify{
			ProcessSearch: func(db *gorm.DB, query string) *gorm.DB {
				return db.Where("first_name || ' ' || last_name like ?", "%"+query+"%")
			},
		})

	if len(req.Type) > 0 {
		db = db.Where("type in (?)", req.Type)
	}

	return sql.GetWithPaginationResponse[entity.Device](db, req.GetManyRequest)
}

func (r DeviceRepository) FindOne(id int64) (entity.Device, error) {
	var device entity.Device
	err := r.db.Where("id = ?", id).First(&device).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Device{}, entity.ErrNotFound
	}

	if err != nil {
		return entity.Device{}, fmt.Errorf("error finding device: %w", err)
	}

	return device, nil
}

func (r DeviceRepository) Create(device *entity.Device) error {
	return r.db.Create(device).Error
}

func (r DeviceRepository) Update(device *entity.Device) error {
	res := r.db.Save(device).Error
	if res != nil {
		return fmt.Errorf("error updating device: %w", res)
	}

	return nil
}

func (r DeviceRepository) Delete(id int) error {
	res := r.db.Delete(&entity.Device{ID: id})
	if res.Error != nil {
		return fmt.Errorf("error deleting device: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return entity.ErrNotFound
	}

	return nil
}

func (r DeviceRepository) FindOneByReceivePort(port string) (entity.Device, error) {
	var device entity.Device
	err := r.db.Where("receive_port = ?", port).First(&device).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Device{}, entity.ErrNotFound
	}

	if err != nil {
		return entity.Device{}, fmt.Errorf("error finding device: %w", err)
	}

	return device, nil
}
