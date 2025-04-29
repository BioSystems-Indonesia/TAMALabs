package entity

import "time"

type SequenceDaily struct {
	ID           int          `gorm:"primaryKey;autoIncrement"`
	SequenceType SequenceType `gorm:"not null"`
	CurrentValue int
	LastUpdated  time.Time
}

// TableName overrides the table name used by DailySequenceModel
func (SequenceDaily) TableName() string {
	return "sequence_daily" // singular table name
}

type SequenceType string

const (
	OrderBarcodeSequence SequenceType = "order"
)
