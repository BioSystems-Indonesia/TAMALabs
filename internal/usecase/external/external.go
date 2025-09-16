package externaluc

import (
	"context"
	"errors"
	"fmt"

	"github.com/oibacidem/lims-hl-seven/config"
	patientrepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/patient"
	testType "github.com/oibacidem/lims-hl-seven/internal/repository/sql/test_type"
	workOrder "github.com/oibacidem/lims-hl-seven/internal/repository/sql/work_order"
	khanzauc "github.com/oibacidem/lims-hl-seven/internal/usecase/external/khanza"
)

type Usecase struct {
	khanzauUC     *khanzauc.Usecase
	workOrderRepo *workOrder.WorkOrderRepository
	patientRepo   *patientrepo.PatientRepository
	testTypeRepo  *testType.Repository
	cfg           *config.Schema
}

func NewUsecase(khanzauUC *khanzauc.Usecase, cfg *config.Schema) *Usecase {
	return &Usecase{
		khanzauUC: khanzauUC,
		cfg:       cfg,
	}
}

func (u *Usecase) SyncAllRequest(ctx context.Context) error {
	var errs []error
	if u.cfg.KhanzaIntegrationEnabled == "true" {
		err := u.khanzauUC.SyncAllRequest(ctx)
		if err != nil {
			errs = append(errs, fmt.Errorf("error syncing all requests khanza: %w", err))
		}
	}

	return errors.Join(errs...)
}

func (u *Usecase) SyncAllResult(ctx context.Context, orderIDs []int64) error {
	var errs []error
	if u.cfg.KhanzaIntegrationEnabled == "true" {
		err := u.khanzauUC.SyncAllResult(ctx, orderIDs)
		if err != nil {
			errs = append(errs, fmt.Errorf("error syncing all results khanza: %w", err))
		}
	}

	return errors.Join(errs...)
}
