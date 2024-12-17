package sql

import "github.com/oibacidem/lims-hl-seven/internal/entity"

type Specimen interface {
	Create(data *entity.Specimen) error
}
