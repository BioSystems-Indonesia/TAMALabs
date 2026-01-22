package usecase

import (
	"context"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
)

type QualityControl interface {
	ParseAndSaveQC(ctx context.Context, hl7Message string, deviceIdentifier string) error
	CreateQCEntry(ctx context.Context, req *entity.CreateQCEntryRequest) (*entity.QCEntry, error)
	CreateManualQCResult(ctx context.Context, req *entity.CreateManualQCResultRequest, createdBy string) (*entity.QCResult, error)
	GetQCEntries(ctx context.Context, req entity.GetManyRequestQCEntry) ([]entity.QCEntry, int64, error)
	GetQCResults(ctx context.Context, req entity.GetManyRequestQCResult) ([]entity.QCResult, int64, error)
	GetQCHistory(ctx context.Context, req entity.GetManyRequestQualityControl) ([]entity.QualityControl, int64, error)
	GetQCStatistics(ctx context.Context, deviceID int) (map[string]interface{}, error)
	UpdateSelectedQCResult(ctx context.Context, qcEntryID int, qcLevel int, resultID int) error
	UpdateSelectedQCResultWithMethod(ctx context.Context, qcEntryID int, qcLevel int, resultID int, method string) error
}
