package admin_uc

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	adminrepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/admin"
	rolerepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/role"
	"golang.org/x/crypto/bcrypt"
)

type AdminUsecase struct {
	adminRepo *adminrepo.AdminRepository
	rolesRepo *rolerepo.RoleRepository
}

func NewAdminUsecase(adminRepo *adminrepo.AdminRepository, roleRepo *rolerepo.RoleRepository) *AdminUsecase {
	return &AdminUsecase{
		adminRepo: adminRepo,
		rolesRepo: roleRepo,
	}
}

func (u *AdminUsecase) CreateAdmin(ctx context.Context, admin *entity.Admin) error {
	if admin.Password == "" {
		return entity.NewHTTPError(http.StatusBadRequest, "password cannot be empty")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin.PasswordHash = string(hashedPassword)

	roleIDs := make([]string, len(admin.RolesID))
	for i, role := range admin.RolesID {
		roleIDs[i] = strconv.Itoa(int(role))
	}

	roles, err := u.rolesRepo.FindAll(ctx, &entity.GetManyRequestRole{
		GetManyRequest: entity.GetManyRequest{
			ID: roleIDs,
		},
	})
	if err != nil {
		return fmt.Errorf("error finding roles: %w", err)
	}
	admin.Roles = roles.Data

	err = u.adminRepo.Create(admin)
	if err != nil {
		return err
	}

	return nil
}

func (u *AdminUsecase) GetAllAdmin(ctx context.Context, req *entity.GetManyRequestAdmin) (entity.PaginationResponse[entity.Admin], error) {
	return u.adminRepo.FindAll(ctx, req)
}

func (u *AdminUsecase) GetOneAdmin(ctx context.Context, id int64) (entity.Admin, error) {
	admin, err := u.adminRepo.FindOne(id)
	if err != nil {
		return entity.Admin{}, err
	}

	roleIDs := make([]int, len(admin.Roles))
	for i, role := range admin.Roles {
		roleIDs[i] = int(role.ID)
	}
	admin.RolesID = roleIDs

	return admin, nil
}

func (u *AdminUsecase) UpdateAdmin(ctx context.Context, admin *entity.Admin) error {
	if admin.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		admin.PasswordHash = string(hashedPassword)
	}

	roleIDs := make([]string, len(admin.RolesID))
	for i, role := range admin.RolesID {
		roleIDs[i] = strconv.Itoa(int(role))
	}

	roles, err := u.rolesRepo.FindAll(ctx, &entity.GetManyRequestRole{
		GetManyRequest: entity.GetManyRequest{
			ID: roleIDs,
		},
	})
	if err != nil {
		return fmt.Errorf("error finding roles: %w", err)
	}
	admin.Roles = roles.Data

	return u.adminRepo.Update(admin)
}

func (u *AdminUsecase) DeleteAdmin(ctx context.Context, id int64) error {
	// Check if admin is still related to any work orders
	err := u.adminRepo.CheckRelatedWorkOrders(ctx, id)
	if err != nil {
		return err
	}

	return u.adminRepo.Delete(id)
}

// GetAllDoctors returns all admins with Doctor role
func (u *AdminUsecase) GetAllDoctors(ctx context.Context) ([]entity.Admin, error) {
	return u.adminRepo.FindAllByRole(ctx, entity.RoleDoctor)
}

// GetAllAnalyzers returns all admins with Analyzer role
func (u *AdminUsecase) GetAllAnalyzers(ctx context.Context) ([]entity.Admin, error) {
	return u.adminRepo.FindAllByRole(ctx, entity.RoleAnalyzer)
}
