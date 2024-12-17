package sql

import "github.com/oibacidem/lims-hl-seven/internal/entity"

type Observation interface {
	Create(data *entity.Observation) error
	CreateMany(data *[]entity.Observation) error
}
