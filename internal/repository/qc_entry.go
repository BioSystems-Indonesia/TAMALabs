package repository

import (
	"context"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
)

type QCEntry interface {
	Create(ctx context.Context, entry *entity.QCEntry) error
	Update(ctx context.Context, id int, entry *entity.UpdateQCEntryRequest) error
	GetByID(ctx context.Context, id int) (*entity.QCEntry, error)
	GetMany(ctx context.Context, req entity.GetManyRequestQCEntry) ([]entity.QCEntry, int64, error)
	GetActiveEntry(ctx context.Context, deviceID, testTypeID, qcLevel int) (*entity.QCEntry, error)
	DeactivateOldEntries(ctx context.Context, deviceID, testTypeID, qcLevel int) error
	GetDeviceSummary(ctx context.Context, deviceID int) (*entity.QCSummary, error)
}

type QCResult interface {
	Create(ctx context.Context, result *entity.QCResult) error
	GetByID(ctx context.Context, id int) (*entity.QCResult, error)
	GetMany(ctx context.Context, req entity.GetManyRequestQCResult) ([]entity.QCResult, int64, error)
	GetByEntryID(ctx context.Context, entryID int) ([]entity.QCResult, error)
	GetByEntryIDAndMethod(ctx context.Context, entryID int, method string) ([]entity.QCResult, error)
	GetByEntryIDAndMethodWithDateRange(ctx context.Context, entryID int, method string, startDate, endDate *string) ([]entity.QCResult, error)
	CalculateStatistics(ctx context.Context, entryID int) (mean float64, sd float64, count int, err error)
	GetCountByLevel(ctx context.Context, deviceID int, testTypeID *int) (map[string]interface{}, error)
}
