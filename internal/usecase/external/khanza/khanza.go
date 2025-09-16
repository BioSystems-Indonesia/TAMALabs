package khanzauc

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

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

	var errs []error
	for _, group := range groupedByOrder {
		firstOrder := group[0]

		patient := groupedPatientBySIMRS[firstOrder.PID]
		if patient.ID == 0 {
			errs = append(errs, fmt.Errorf("patient not found with PID: %s", firstOrder.PID))
			continue
		}

		err = u.insertOrderToLIS(ctx, firstOrder.ONO, firstOrder.VisitNo, patient)
		if err != nil {
			errs = append(errs, fmt.Errorf("error syncing request: %w", err))
			continue
		}
	}

	return errors.Join(errs...)
}

func (u *Usecase) insertOrderToLIS(
	ctx context.Context,
	ono string,
	visitNo string,
	patient entity.Patient,
) error {
	lisOrders, err := u.repo.GetLabRequestByNoOrder(ctx, ono)
	if err != nil {
		return fmt.Errorf("error getting tests by ONO: %w", err)
	}

	var tests []entity.WorkOrderCreateRequestTestType
	for _, lisOrder := range lisOrders {
		testType, err := u.testTypeRepo.FindOneByAliasCode(ctx, lisOrder.Pemeriksaan)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				slog.WarnContext(ctx, "khanza test type not found", "aliasCode", lisOrder.Pemeriksaan)
				continue
			}

			slog.WarnContext(ctx, "error getting test type by alias code", "aliasCode", lisOrder.Pemeriksaan, "error", err)
			continue
		}

		tests = append(tests, entity.WorkOrderCreateRequestTestType{
			TestTypeID:   int64(testType.ID),
			TestTypeCode: testType.Code,
			SpecimenType: testType.GetFirstType(),
		})
	}

	barcode, err := u.barcodeUC.NextOrderBarcode(ctx)
	if err != nil {
		return fmt.Errorf("failed to u.barcodeGeneratorUC.NextOrderBarcode %w", err)
	}

	createReq := entity.WorkOrderCreateRequest{
		PatientID:    patient.ID,
		Barcode:      barcode,
		BarcodeSIMRS: visitNo,
		TestTypes:    tests,
		CreatedBy:    -1,
	}

	existingWorkOrder, err := u.workOrderRepository.GetBySIMRSBarcode(ctx, visitNo)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			_, err = u.workOrderRepository.Create(&createReq)
			if err != nil {
				return fmt.Errorf("error creating work order: %w", err)
			}
			return nil
		}

		return fmt.Errorf("error checking work order by barcode: %w", err)
	}

	if existingWorkOrder.Status != entity.WorkOrderStatusNew {
		return nil
	}

	_, err = u.workOrderRepository.Edit(int(existingWorkOrder.ID), &createReq)
	if err != nil {
		return fmt.Errorf("error updating work order: %w", err)
	}
	return nil
}

func (u *Usecase) ProcessRequest(ctx context.Context, rawRequest []byte) error {
	var request Request
	var err error

	slog.Info("debug khanza process request", "rawRequest", string(rawRequest))

	err = json.Unmarshal(rawRequest, &request)
	if err != nil {
		return fmt.Errorf("external.ProcessRequest json.Unmarshal failed: %w\nbody: %s", err, string(rawRequest))
	}

	var patient entity.Patient

	patient, err = u.findOrCreatePatient(request)
	if err != nil {
		return fmt.Errorf("external.ProcessRequest findOrCreatePatient failed: %w", err)
	}

	if err := u.insertOrderToLIS(ctx, request.Order.OBR.OrderLab, request.Order.OBR.RegNo, patient); err != nil {
		return fmt.Errorf("external.ProcessRequest insertOrderToLIS failed: %w", err)
	}

	return nil
}

func (u *Usecase) GetResult(ctx context.Context, ono string) ([]byte, error) {
	res := Response{}

	slog.InfoContext(ctx, "debug khanza get result", "ono", ono)

	orders, err := u.repo.GetLabRequestByNoOrder(ctx, ono)
	if err != nil {
		return nil, fmt.Errorf("error getting order by no order: %w", err)
	}

	if len(orders) == 0 {
		return nil, fmt.Errorf("order not found")
	}

	firstOrders := orders[0]
	workOrder, err := u.workOrderRepository.GetBySIMRSBarcode(ctx, firstOrders.NoRawat)
	if err != nil {
		return nil, err
	}
	workOrder.FillTestResultDetail(false)

	res.Result.OBX.OrderLab = ono

	resultTest := make([]ResponseResultTest, len(workOrder.TestResult))

	for i, t := range workOrder.TestResult {
		hasil := ""
		if t.Result != nil {
			hasil = strconv.FormatFloat(*t.Result, 'f', 4, 64)
		}

		testName := t.TestType.AliasCode
		if testName == "" {
			testName = t.TestType.Name
		}

		resultTest[i] = ResponseResultTest{
			TestID:      strconv.FormatInt(t.ID, 10),
			NamaTest:    testName,
			Hasil:       hasil,
			NilaiNormal: t.ReferenceRange,
			Satuan:      t.Unit,
			Flag:        "",
		}
	}

	resByte, err := json.Marshal(res)
	if err != nil {
		return nil, fmt.Errorf("error marshaling response: %w", err)
	}

	slog.InfoContext(ctx, "debug khanza get result", "res", string(resByte))

	return resByte, nil
}

func (u *Usecase) findOrCreatePatient(r Request) (entity.Patient, error) {
	p, err := u.convertIntoPatient(r)
	if err != nil {
		return p, fmt.Errorf("error converting patient: %w", err)
	}

	patients := []entity.Patient{p}
	err = u.patientRepository.CreateManyFromSIMRS(patients)
	if err != nil {
		return p, fmt.Errorf("error creating patient: %w", err)
	}

	p = patients[0]
	if p.ID == 0 {
		return p, errors.New("patient id is zero, id is not filled on DB")
	}

	return p, nil
}

func (u *Usecase) convertIntoPatient(r Request) (entity.Patient, error) {
	firstName, lastName := util.SplitName(r.Order.PID.Pname)

	sex := entity.PatientSexUnknown
	switch r.Order.PID.Sex {
	case "L":
		sex = entity.PatientSexMale
	case "F":
		sex = entity.PatientSexFemale
	}

	birthdate, err := time.Parse("02.01.2006", r.Order.PID.BirthDt)
	if err != nil {
		return entity.Patient{}, fmt.Errorf("cannot parse birth date: %w", err)
	}

	phoneNumber := r.Order.PID.NoHP
	if phoneNumber == "" {
		phoneNumber = r.Order.PID.NoTlp
	}

	return entity.Patient{
		FirstName: firstName,
		LastName:  lastName,
		Sex:       sex,
		Birthdate: birthdate,
		SIMRSPID: sql.NullString{
			String: r.Order.PID.Pmrn,
			Valid:  true,
		},
		Address:     r.Order.PID.Address,
		PhoneNumber: phoneNumber,
		Location:    u.getLocation(r),
	}, nil
}

func (u *Usecase) getLocation(r Request) string {
	bangsalLocation := []string{}
	if r.Order.OBR.BangsalID != "" {
		bangsalLocation = append(bangsalLocation, r.Order.OBR.BangsalID)
	}
	if r.Order.OBR.BedID != "" {
		bangsalLocation = append(bangsalLocation, r.Order.OBR.BangsalName)
	}

	bedLocation := []string{}
	if r.Order.OBR.BedID != "" {
		bedLocation = append(bedLocation, r.Order.OBR.BedID)
	}
	if r.Order.OBR.BedName != "" {
		bedLocation = append(bedLocation, r.Order.OBR.BedName)
	}

	allLocation := []string{}
	allLocation = append(allLocation, strings.Join(bangsalLocation, "-"))
	allLocation = append(allLocation, strings.Join(bedLocation, "-"))

	return strings.Join(allLocation, "|")
}
