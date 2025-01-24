package configuc

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	configrepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/config"
)

type ConfigUseCase struct {
	cfg        *config.Schema
	configRepo *configrepo.Repository
	validate   *validator.Validate
}

func NewConfigUseCase(
	cfg *config.Schema,
	ConfigRepo *configrepo.Repository,
	validate *validator.Validate,
) *ConfigUseCase {
	return &ConfigUseCase{cfg: cfg, configRepo: ConfigRepo, validate: validate}
}

func (p ConfigUseCase) FindAll(
	ctx context.Context, req *entity.ConfigGetManyRequest,
) (entity.PaginationResponse[entity.Config], error) {
	return p.configRepo.FindAll(ctx, req)
}

func (p ConfigUseCase) FindOneByID(ctx context.Context, key string) (entity.Config, error) {
	return p.configRepo.FindOne(ctx, key)
}

func (p ConfigUseCase) Edit(ctx context.Context, key string, value string) (entity.Config, error) {
	config, err := p.configRepo.Edit(ctx, key, value)
	if err != nil {
		return entity.Config{}, err
	}
	return config, nil
}
