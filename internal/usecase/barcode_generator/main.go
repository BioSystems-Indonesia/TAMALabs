package barcode_generator

import (
	"context"
	"fmt"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/daily_sequence"
)

type Usecase struct {
	dailySequenceRepo *daily_sequence.Repository
}

func NewUsecase(dailySequence *daily_sequence.Repository) *Usecase {
	return &Usecase{dailySequenceRepo: dailySequence}
}

func (u *Usecase) NextOrderBarcode(ctx context.Context) (string, error) {
	now := time.Now()

	// Use atomic GetNextSequence instead of separate GetOrReset + Incr
	nextSeq, err := u.dailySequenceRepo.GetNextSequence(ctx, now, entity.OrderBarcodeSequence)
	if err != nil {
		return "", fmt.Errorf("failed to get next sequence: %w", err)
	}

	return fmt.Sprintf("%s%03d", now.Format("060102"), nextSeq), nil
}
