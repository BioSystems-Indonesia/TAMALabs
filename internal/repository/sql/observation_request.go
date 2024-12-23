package sql

import "github.com/oibacidem/lims-hl-seven/internal/entity"

type ObservationRequest interface {
	Create(data *entity.ObservationRequest) error
	CreateMany(data *[]entity.ObservationRequest) error
}
