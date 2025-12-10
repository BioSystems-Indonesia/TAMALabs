package adminrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/BioSystems-Indonesia/TAMALabs/config"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql"
	"gorm.io/gorm"
)

type AdminRepository struct {
	db  *gorm.DB
	cfg *config.Schema
}

func NewAdminRepository(db *gorm.DB, cfg *config.Schema) *AdminRepository {
	return &AdminRepository{db: db, cfg: cfg}
}

func (r AdminRepository) FindAll(
	ctx context.Context,
	req *entity.GetManyRequestAdmin,
) (entity.PaginationResponse[entity.Admin], error) {
	var roleIDs []int64
	if len(req.Role) > 0 {
		dbRole := r.db.WithContext(ctx).Model(&entity.Role{})
		dbRole = dbRole.Where("name in ?", req.Role)
		err := dbRole.Pluck("id", &roleIDs).Error
		if err != nil {
			return entity.PaginationResponse[entity.Admin]{}, fmt.Errorf("error finding roles: %w", err)
		}
	}

	db := r.db.WithContext(ctx).Preload("Roles")
	db = sql.ProcessGetMany(db, req.GetManyRequest,
		sql.Modify{
			ProcessSearch: func(db *gorm.DB, query string) *gorm.DB {
				db = db.Where("fullname like ? or email like ? or username like ?", query+"%", query+"%", query+"%")
				return db
			},
		})

	if len(roleIDs) > 0 {
		db = db.Distinct()
		db = db.Joins("join admin_roles on admin_roles.admin_id = admins.id and admin_roles.role_id in?", roleIDs)
	}

	return sql.GetWithPaginationResponse[entity.Admin](db, req.GetManyRequest)
}

// FindAllByRole finds all admins with a specific role without pagination
func (r AdminRepository) FindAllByRole(ctx context.Context, roleName entity.RoleName) ([]entity.Admin, error) {
	var admins []entity.Admin

	err := r.db.WithContext(ctx).
		Preload("Roles").
		Joins("JOIN admin_roles ON admin_roles.admin_id = admins.id").
		Joins("JOIN roles ON roles.id = admin_roles.role_id AND roles.name = ?", roleName).
		Where("admins.is_active = ?", true).
		Order("admins.fullname ASC").
		Find(&admins).Error

	if err != nil {
		return nil, fmt.Errorf("error finding admins by role: %w", err)
	}

	return admins, nil
}

func (r AdminRepository) FindOne(id int64) (entity.Admin, error) {
	var admin entity.Admin
	err := r.db.Where("id = ?", id).Preload("Roles").First(&admin).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Admin{}, entity.ErrNotFound
	}

	if err != nil {
		return entity.Admin{}, fmt.Errorf("error finding admin: %w", err)
	}

	return admin, nil
}

func (r AdminRepository) Create(admin *entity.Admin) error {
	return r.db.Create(admin).Error
}

func (r AdminRepository) Update(admin *entity.Admin) error {
	res := r.db.Updates(admin).Error
	if res != nil {
		return fmt.Errorf("error updating admin: %w", res)
	}

	return nil
}

func (r AdminRepository) Delete(id int64) error {
	res := r.db.Delete(&entity.Admin{ID: id})
	if res.Error != nil {
		return fmt.Errorf("error deleting admin: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return entity.ErrNotFound
	}

	return nil
}

func (r AdminRepository) CheckRelatedWorkOrders(ctx context.Context, adminID int64) error {
	// Check if admin is still related to any work orders as doctor
	var doctorCount int64
	err := r.db.WithContext(ctx).Table("work_order_doctors").Where("admin_id = ?", adminID).Count(&doctorCount).Error
	if err != nil {
		return fmt.Errorf("error checking work order doctors relation: %w", err)
	}

	if doctorCount > 0 {
		return entity.NewHTTPError(
			entity.ErrCannotDeleteAdminWithRelations.Code,
			fmt.Sprintf("Tidak dapat menghapus admin: masih ditugaskan sebagai dokter di %d work order", doctorCount),
		)
	}

	// Check if admin is still related to any work orders as analyzer
	var analyzerCount int64
	err = r.db.WithContext(ctx).Table("work_order_analyzers").Where("admin_id = ?", adminID).Count(&analyzerCount).Error
	if err != nil {
		return fmt.Errorf("error checking work order analyzers relation: %w", err)
	}

	if analyzerCount > 0 {
		return entity.NewHTTPError(
			entity.ErrCannotDeleteAdminWithRelations.Code,
			fmt.Sprintf("Tidak dapat menghapus admin: masih ditugaskan sebagai analyzer di %d work order", analyzerCount),
		)
	}

	// Check if admin created any work orders
	var createdOrdersCount int64
	err = r.db.WithContext(ctx).Table("work_orders").Where("created_by = ?", adminID).Count(&createdOrdersCount).Error
	if err != nil {
		return fmt.Errorf("error checking created work orders: %w", err)
	}

	if createdOrdersCount > 0 {
		return entity.NewHTTPError(
			entity.ErrCannotDeleteAdminWithRelations.Code,
			fmt.Sprintf("Tidak dapat menghapus admin: masih tercatat sebagai pembuat di %d work order", createdOrdersCount),
		)
	}

	// Check if admin last updated any work orders
	var updatedOrdersCount int64
	err = r.db.WithContext(ctx).Table("work_orders").Where("last_updated_by = ?", adminID).Count(&updatedOrdersCount).Error
	if err != nil {
		return fmt.Errorf("error checking updated work orders: %w", err)
	}

	if updatedOrdersCount > 0 {
		return entity.NewHTTPError(
			entity.ErrCannotDeleteAdminWithRelations.Code,
			fmt.Sprintf("Tidak dapat menghapus admin: masih tercatat sebagai yang terakhir mengupdate di %d work order", updatedOrdersCount),
		)
	}

	// Check if admin created any test templates
	var testTemplateCreatedCount int64
	err = r.db.WithContext(ctx).Table("test_templates").Where("created_by = ?", adminID).Count(&testTemplateCreatedCount).Error
	if err != nil {
		return fmt.Errorf("error checking test template created by admin: %w", err)
	}

	if testTemplateCreatedCount > 0 {
		return entity.NewHTTPError(
			entity.ErrCannotDeleteAdminWithRelations.Code,
			fmt.Sprintf("Tidak dapat menghapus admin: masih tercatat sebagai pembuat di %d test template", testTemplateCreatedCount),
		)
	}

	// Check if admin last updated any test templates
	var testTemplateUpdatedCount int64
	err = r.db.WithContext(ctx).Table("test_templates").Where("last_updated_by = ?", adminID).Count(&testTemplateUpdatedCount).Error
	if err != nil {
		return fmt.Errorf("error checking test template updated by admin: %w", err)
	}

	if testTemplateUpdatedCount > 0 {
		return entity.NewHTTPError(
			entity.ErrCannotDeleteAdminWithRelations.Code,
			fmt.Sprintf("Tidak dapat menghapus admin: masih tercatat sebagai yang terakhir mengupdate di %d test template", testTemplateUpdatedCount),
		)
	}

	return nil
}

func (r AdminRepository) FindOneByEmail(ctx context.Context, email string) (entity.Admin, error) {
	var admin entity.Admin
	err := r.db.Where("email = ?", email).First(&admin).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Admin{}, entity.ErrNotFound
	}

	if err != nil {
		return entity.Admin{}, fmt.Errorf("error finding admin: %w", err)
	}

	return admin, nil
}

func (r AdminRepository) FindOneByUsername(ctx context.Context, username string) (entity.Admin, error) {
	var admin entity.Admin
	err := r.db.Where("username = ?", username).Preload("Roles").First(&admin).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Admin{}, entity.ErrNotFound
	}

	if err != nil {
		return entity.Admin{}, fmt.Errorf("error finding admin: %w", err)
	}

	return admin, nil
}
