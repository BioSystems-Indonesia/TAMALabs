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

	"github.com/labstack/gommon/log"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/repository/external/khanza"
	patientrepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/patient"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/test_type"
	workOrderRepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/work_order"
	"github.com/oibacidem/lims-hl-seven/internal/usecase"
	"github.com/oibacidem/lims-hl-seven/internal/usecase/result"
	"github.com/oibacidem/lims-hl-seven/internal/util"
	"gorm.io/gorm"
)

type Usecase struct {
	repo                *khanza.Repository
	patientRepository   *patientrepo.PatientRepository
	testTypeRepo        *test_type.Repository
	workOrderRepository *workOrderRepo.WorkOrderRepository
	barcodeUC           usecase.BarcodeGenerator
	resultUC            *result.Usecase
}

func NewUsecase(
	repo *khanza.Repository,
	workOrderRepository *workOrderRepo.WorkOrderRepository,
	patientRepository *patientrepo.PatientRepository,
	testTypeRepo *test_type.Repository,
	barcodeUC usecase.BarcodeGenerator,
	resultUC *result.Usecase,
) *Usecase {
	return &Usecase{
		repo:                repo,
		patientRepository:   patientRepository,
		testTypeRepo:        testTypeRepo,
		workOrderRepository: workOrderRepository,
		barcodeUC:           barcodeUC,
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
		slog.Info("testResult", "testResult", testResult)
		barcode := workOrder.BarcodeSIMRS // visit no
		if barcode == "" {
			continue
		}

		// Translate visit no to ono
		order, err := u.repo.GetLisOrderByVisitNo(barcode)
		if err != nil {
			return fmt.Errorf("error getting lis order by visit no: %w", err)
		}
		barcode = order.ONO

		code := testResult.TestType.AliasCode
		if code == "" {
			code = testResult.TestType.Code
		}

		khanzaDT := entity.KhanzaResDT{
			ONO:         barcode,
			TESTCD:      code,
			TestNM:      code,
			DataTyp:     entity.DataTypNumeric,
			ResultValue: fmt.Sprintf("%.2f", testResult.GetResult()),
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
			Birthdate: firstOrder.BirthDT,
			Sex:       entity.NewPatientSexFromKhanza(entity.KhanzaPatientSex(firstOrder.Sex)),
			Address:   address,
		}
		patients = append(patients, patient)
	}

	err = u.patientRepository.CreateManyFromSIMRS(patients)
	if err != nil {
		return fmt.Errorf("error creating patient: %w", err)
	}

	groupedPatientBySIMRS := make(map[string]entity.Patient)
	for _, patient := range patients {
		groupedPatientBySIMRS[patient.SIMRSPID.String] = patient
	}

	groupedByOrder := make(map[string][]entity.KhanzaLisOrder)
	for _, order := range orders {
		groupedByOrder[order.ONO] = append(groupedByOrder[order.ONO], order)
	}

	for _, group := range groupedByOrder {
		firstOrder := group[0]

		patient := groupedPatientBySIMRS[firstOrder.PID]
		if patient.ID == 0 {
			return fmt.Errorf("patient not found with PID: %s", firstOrder.PID)
		}

		lisOrders, err := u.repo.GetLabRequestByNoOrder(ctx, firstOrder.ONO)
		if err != nil {
			return fmt.Errorf("error getting tests by ONO: %w", err)
		}

		var tests []entity.WorkOrderCreateRequestTestType
		for _, lisOrder := range lisOrders {
			testType, err := u.testTypeRepo.FindOneByAliasCode(ctx, lisOrder.Pemeriksaan)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					log.Warnf("test type not found with alias code: %s", lisOrder.Pemeriksaan)
					continue
				}

				log.Warnf("error getting test type by alias code: %s", lisOrder.Pemeriksaan)
				continue
			}

			tests = append(tests, entity.WorkOrderCreateRequestTestType{
				TestTypeID:   int64(testType.ID),
				TestTypeCode: testType.Code,
				SpecimenType: testType.GetFirstType(),
			})
		}

		barcode, err := u.barcodeUC.NextOrderBarcode(context.Background())
		if err != nil {
			return fmt.Errorf("failed to u.barcodeGeneratorUC.NextOrderBarcode %w", err)
		}

		createReq := entity.WorkOrderCreateRequest{
			PatientID:    patient.ID,
			Barcode:      barcode,
			BarcodeSIMRS: firstOrder.VisitNo,
			TestTypes:    tests,
			CreatedBy:    -1,
		}

		existingWorkOrder, err := u.workOrderRepository.GetBySIMRSBarcode(ctx, firstOrder.VisitNo)
		if err != nil {
			if errors.Is(err, entity.ErrNotFound) {
				_, err = u.workOrderRepository.Create(&createReq)
				if err != nil {
					return fmt.Errorf("error creating work order: %w", err)
				}
				continue
			}

			return fmt.Errorf("error checking work order by barcode: %w", err)
		}

		if existingWorkOrder.Status != entity.WorkOrderStatusNew {
			continue
		}

		_, err = u.workOrderRepository.Edit(int(existingWorkOrder.ID), &createReq)
		if err != nil {
			return fmt.Errorf("error updating work order: %w", err)
		}
	}

	return nil
}
