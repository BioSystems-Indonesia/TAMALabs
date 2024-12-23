package sql

import "github.com/oibacidem/lims-hl-seven/internal/entity"

type ObservationResult interface {
	Create(data *entity.ObservationResult) error
	CreateMany(data *[]entity.ObservationResult) error
}
