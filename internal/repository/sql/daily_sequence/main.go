package daily_sequence

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
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
	var sequence entity.SequenceDaily

	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Add timeout for the query to prevent hanging - increased for midnight operations
		queryCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
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

		// Check if we need to reset for new day using date truncation for accuracy
		currentDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
		lastUpdatedDay := time.Date(sequence.LastUpdated.Year(), sequence.LastUpdated.Month(), sequence.LastUpdated.Day(), 0, 0, 0, 0, time.Local)

		if !currentDay.Equal(lastUpdatedDay) {
			sequence.CurrentValue = 0
			sequence.LastUpdated = currentDay
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
	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Add timeout for increment operation
		queryCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		if err := tx.WithContext(queryCtx).Clauses(
			clause.Locking{Strength: "UPDATE"},
		).First(&sequence, "sequence_type = ?", seqType).Error; err != nil {
			return fmt.Errorf("sequence not found: %w", err)
		}

		sequence.CurrentValue++
		sequence.LastUpdated = time.Now()
		if err := tx.WithContext(queryCtx).Save(&sequence).Error; err != nil {
			return fmt.Errorf("failed to increment sequence: %w", err)
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return sequence.CurrentValue, nil
}

// GetNextSequence atomically gets the next sequence number with proper daily reset
func (r *Repository) GetNextSequence(ctx context.Context, date time.Time, seqType entity.SequenceType) (int, error) {
	var nextValue int

	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		queryCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		var sequence entity.SequenceDaily
		result := tx.WithContext(queryCtx).Clauses(
			clause.Locking{Strength: "UPDATE"},
		).First(&sequence, "sequence_type = ?", seqType)

		if result.Error == gorm.ErrRecordNotFound {
			// First time, create with value 1
			sequence = entity.SequenceDaily{
				SequenceType: seqType,
				CurrentValue: 1,
				LastUpdated:  date,
			}
			nextValue = 1

			if err := tx.WithContext(queryCtx).Create(&sequence).Error; err != nil {
				return fmt.Errorf("failed to create sequence: %w", err)
			}
			return nil
		}

		if result.Error != nil {
			return fmt.Errorf("failed to get sequence: %w", result.Error)
		}

		// Check if we need to reset for new day using date truncation for accuracy
		currentDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
		lastUpdatedDay := time.Date(sequence.LastUpdated.Year(), sequence.LastUpdated.Month(), sequence.LastUpdated.Day(), 0, 0, 0, 0, time.Local)

		if !currentDay.Equal(lastUpdatedDay) {
			// New day, reset to 1
			sequence.CurrentValue = 1
			sequence.LastUpdated = currentDay
			nextValue = 1
		} else {
			// Same day, increment
			sequence.CurrentValue++
			nextValue = sequence.CurrentValue
		}

		if err := tx.WithContext(queryCtx).Save(&sequence).Error; err != nil {
			return fmt.Errorf("failed to update sequence: %w", err)
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return nextValue, nil
}
