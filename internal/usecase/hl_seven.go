package usecase

import "github.com/oibacidem/lims-hl-seven/internal/entity"

// HLSeven is an interface for HLSeven usecase
type HLSeven interface {
	SendORM(message entity.SendORMRequest) (*entity.SendORMResponse, error)
	ProcessORM(orm entity.ORM) (string, error)
}
