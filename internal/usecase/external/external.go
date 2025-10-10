package externaluc

import (
	"context"
	"errors"
	"fmt"

	"github.com/BioSystems-Indonesia/TAMALabs/config"
	simrs "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/external/simrs"
	patientrepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/patient"
	testType "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/test_type"
	workOrder "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/work_order"
	khanzauc "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/external/khanza"
	simrsuc "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/external/simrs"
)

type Usecase struct {
	khanzauUC     *khanzauc.Usecase
	simrsUC       *simrsuc.Usecase
	workOrderRepo *workOrder.WorkOrderRepository
	patientRepo   *patientrepo.PatientRepository
	testTypeRepo  *testType.Repository
	cfg           *config.Schema
}

func NewUsecase(khanzauUC *khanzauc.Usecase, simrsUC *simrsuc.Usecase, workOrderRepo *workOrder.WorkOrderRepository, cfg *config.Schema) *Usecase {
	return &Usecase{
		khanzauUC:     khanzauUC,
		simrsUC:       simrsUC,
		workOrderRepo: workOrderRepo,
		cfg:           cfg,
	}
}

func (u *Usecase) SyncAllRequest(ctx context.Context) error {
	var errs []error

	if u.cfg.KhanzaIntegrationEnabled == "true" && u.khanzauUC != nil {
		err := u.khanzauUC.SyncAllRequest(ctx)
		if err != nil {
			errs = append(errs, fmt.Errorf("error syncing all requests khanza: %w", err))
		}
	}

	if u.cfg.SimrsIntegrationEnabled == "true" && u.simrsUC != nil {
		err := u.simrsUC.SyncAllRequest(ctx)
		if err != nil {
			errs = append(errs, fmt.Errorf("error syncing all requests simrs: %w", err))
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

	if u.cfg.SimrsIntegrationEnabled == "true" && u.simrsUC != nil {
		err := u.simrsUC.SyncAllResult(ctx, orderIDs)
		if err != nil {
			errs = append(errs, fmt.Errorf("error syncing all results simrs: %w", err))
		}
	}

	return errors.Join(errs...)
}

// TestSimrsConnection tests SIMRS database connection
func (u *Usecase) TestSimrsConnection(ctx context.Context, dsn string) error {
	// Try to create a connection with the provided DSN
	simrsDB, err := simrs.NewDB(dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to SIMRS database: %w", err)
	}
	defer simrsDB.Close()

	// Test the connection by pinging
	if err := simrsDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping SIMRS database: %w", err)
	}

	return nil
}
