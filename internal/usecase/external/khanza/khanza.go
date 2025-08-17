package khanzauc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/repository/external/khanza"
	patientrepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/patient"
	workOrderRepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/work_order"
	"github.com/oibacidem/lims-hl-seven/internal/usecase/result"
	"github.com/oibacidem/lims-hl-seven/internal/util"
)

type Usecase struct {
	repo                *khanza.Repository
	patientRepository   *patientrepo.PatientRepository
	workOrderRepository *workOrderRepo.WorkOrderRepository
	resultUC            *result.Usecase
}

func NewUsecase(
	repo *khanza.Repository,
	workOrderRepository *workOrderRepo.WorkOrderRepository,
	patientRepository *patientrepo.PatientRepository,
	resultUC *result.Usecase,
) *Usecase {
	return &Usecase{
		repo:                repo,
		patientRepository:   patientRepository,
		workOrderRepository: workOrderRepository,
		resultUC:            resultUC,
	}
}

func (u *Usecase) SyncAllResult(ctx context.Context, orderIDs []int64) error {
	orderIDStrings := make([]string, len(orderIDs))
	for i, orderID := range orderIDs {
		orderIDStrings[i] = strconv.Itoa(int(orderID))
	}

	workOrders, err := u.workOrderRepository.FindAll(ctx, &entity.WorkOrderGetManyRequest{
		GetManyRequest: entity.GetManyRequest{
			CreatedAtStart: time.Now().Add(14 * -24 * time.Hour),
			CreatedAtEnd:   time.Now(),
			ID:             orderIDStrings,
		},
	})
	if err != nil {
		return err
	}

	var errs []error
	for _, workOrder := range workOrders.Data {
		err := u.SyncResult(ctx, workOrder.ID)
		if err != nil {
			slog.Error("error syncing result", "error", err)
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (u *Usecase) SyncResult(ctx context.Context, workOrderID int64) error {
	workOrder, err := u.resultUC.ResultDetail(ctx, workOrderID)
	if err != nil {
		return err
	}

	var tests []entity.TestResult
	for _, testResult := range workOrder.TestResult {
		tests = append(tests, testResult...)
	}

	var reqs []entity.KhanzaResDT
	for _, testResult := range tests {
		barcode := workOrder.BarcodeSIMRS
		if barcode == "" {
			barcode = workOrder.Barcode
		}

		code := testResult.TestType.AliasCode
		if code == "" {
			code = testResult.TestType.Code
		}

		khanzaDT := entity.KhanzaResDT{
			ONO:         barcode,
			TESTCD:      code,
			TestNM:      code,
			DataTyp:     entity.DataTypNumeric,
			ResultValue: fmt.Sprintf("%.2f", testResult.GetFormattedResult()),
			Unit:        testResult.Unit,
			Flag:        entity.NewKhanzaFlag(testResult),
			RefRange:    testResult.ReferenceRange,
		}
		reqs = append(reqs, khanzaDT)
	}

	err = u.repo.BatchUpsertRESDTO(ctx, reqs)
	if err != nil {
		return fmt.Errorf("error batch upsert: %w", err)
	}

	return nil
}

func (u *Usecase) SyncAllRequest(ctx context.Context) error {
	orders, err := u.repo.GetAllLisOrders()
	if err != nil {
		return err
	}

	// Insert patient PID
	groupedByPID := make(map[string][]entity.KhanzaLisOrder)
	for _, order := range orders {
		groupedByPID[order.PID] = append(groupedByPID[order.PID], order)
	}

	slog.InfoContext(ctx, "grouped by PID", "count", len(groupedByPID))

	var patients []entity.Patient
	for _, group := range groupedByPID {
		firstOrder := group[0]
		firstName, lastName := util.SplitName(firstOrder.PName)
		addresses := []string{
			firstOrder.Address1,
			firstOrder.Address2,
			firstOrder.Address3,
			firstOrder.Address4,
		}
		address := strings.Join(addresses, ",")

		patient := entity.Patient{
			SIMRSPID: sql.NullString{
				String: firstOrder.PID,
				Valid:  true,
			},
			FirstName: firstName,
			LastName:  lastName,
			Birthdate: time.Time{},
			Sex:       entity.NewPatientSexFromKhanza(entity.KhanzaPatientSex(firstOrder.Sex)),
			Address:   address,
		}
		patients = append(patients, patient)
	}

	return u.patientRepository.CreateManyFromSIMRS(patients)
}
