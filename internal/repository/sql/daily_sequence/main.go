package daily_sequence

import (
	"context"
	"fmt"
	"time"

	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) GetOrReset(ctx context.Context, date time.Time, seqType entity.SequenceType) (int, error) {
	var sequence entity.SequenceDaily

	err := r.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Clauses(
			clause.Locking{Strength: "UPDATE"},
		).First(&sequence, "sequence_type = ?", seqType)

		if result.Error == gorm.ErrRecordNotFound {
			sequence = entity.SequenceDaily{
				SequenceType: seqType,
				CurrentValue: 0,
				LastUpdated:  date,
			}

			if err := tx.Create(&sequence).Error; err != nil {
				return fmt.Errorf("failed to create sequence: %w", err)
			}
			return nil
		}

		if result.Error != nil {
			return fmt.Errorf("failed to get sequence: %w", result.Error)
		}
		if sequence.LastUpdated.Format("2006-01-02") != date.Format("2006-01-02") {
			sequence.CurrentValue = 0
			sequence.LastUpdated = date
			if err := tx.Save(&sequence).Error; err != nil {
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
