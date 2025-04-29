package barcode_generator

import (
	"context"
	"fmt"
	"time"

	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/daily_sequence"
)

type Usecase struct {
	dailySequenceRepo *daily_sequence.Repository
}

func NewUsecase(dailySequence *daily_sequence.Repository) *Usecase {
	return &Usecase{dailySequenceRepo: dailySequence}
}

func (u *Usecase) NextOrderBarcode(ctx context.Context) (string, error) {
	now := time.Now()
	seq, err := u.dailySequenceRepo.GetOrReset(ctx, now, entity.OrderBarcodeSequence)
	if err != nil {
		return "", fmt.Errorf("failed to u.dailySequenceRepo.GetOrReset: %w", err)
	}

	nextSeq, err := u.dailySequenceRepo.Incr(ctx, entity.OrderBarcodeSequence, seq)
	if err != nil {
		return "", fmt.Errorf("failed to u.dailySequenceRepo.Incr: %w", err)
	}

	return fmt.Sprintf("%s%s", now.Format("20060102"), fmt.Sprintf("%03d", nextSeq)), nil
}
