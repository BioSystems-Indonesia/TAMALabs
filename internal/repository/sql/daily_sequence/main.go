package daily_sequence

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	DB *gorm.DB
	mu sync.RWMutex
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) GetOrReset(ctx context.Context, date time.Time, seqType entity.SequenceType) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var sequence entity.SequenceDaily

	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Add timeout for the query to prevent hanging
		queryCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		result := tx.WithContext(queryCtx).Clauses(
			clause.Locking{Strength: "UPDATE"},
		).First(&sequence, "sequence_type = ?", seqType)

		if result.Error == gorm.ErrRecordNotFound {
			sequence = entity.SequenceDaily{
				SequenceType: seqType,
				CurrentValue: 0,
				LastUpdated:  date,
			}

			if err := tx.WithContext(queryCtx).Create(&sequence).Error; err != nil {
				return fmt.Errorf("failed to create sequence: %w", err)
			}
			return nil
		}

		if result.Error != nil {
			return fmt.Errorf("failed to get sequence: %w", result.Error)
		}

		// Check if we need to reset for new day
		if sequence.LastUpdated.Format("2006-01-02") != date.Format("2006-01-02") {
			sequence.CurrentValue = 0
			sequence.LastUpdated = date
			if err := tx.WithContext(queryCtx).Save(&sequence).Error; err != nil {
				return fmt.Errorf("failed to reset sequence: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return sequence.CurrentValue, nil
}

func (r *Repository) Incr(ctx context.Context, seqType entity.SequenceType, currentValue int) (int, error) {
	var sequence entity.SequenceDaily
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(
			clause.Locking{Strength: "UPDATE"},
		).First(&sequence, "sequence_type = ?", seqType).Error; err != nil {
			return fmt.Errorf("sequence not found: %w", err)
		}

		sequence.CurrentValue++
		sequence.LastUpdated = time.Now()
		if err := tx.Save(&sequence).Error; err != nil {
			return fmt.Errorf("failed to increment sequence: %w", err)
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return sequence.CurrentValue, nil
}
