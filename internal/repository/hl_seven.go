package repository

import "github.com/oibacidem/lims-hl-seven/internal/entity"

type HLSeven interface {
	SendORM(request entity.SendORMRequest) (*entity.SendORMResponse, error)
}
