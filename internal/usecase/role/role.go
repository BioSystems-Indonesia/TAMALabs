package role_uc

import (
	"context"

	"github.com/oibacidem/lims-hl-seven/internal/entity"
	rolerepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/role"
)

type RoleUsecase struct {
	roleRepo *rolerepo.RoleRepository
}

func NewRoleUsecase(roleRepo *rolerepo.RoleRepository) *RoleUsecase {
	return &RoleUsecase{
		roleRepo: roleRepo,
	}
}

func (u *RoleUsecase) GetAllRole(ctx context.Context, req *entity.GetManyRequestRole) (entity.PaginationResponse[entity.Role], error) {
	return u.roleRepo.FindAll(ctx, req)
}

func (u *RoleUsecase) GetOneRole(ctx context.Context, id int64) (entity.Role, error) {
	return u.roleRepo.FindOne(id)
}
