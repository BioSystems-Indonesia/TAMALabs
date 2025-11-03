package khanzauc

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/repository/external/khanza"
	patientrepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/patient"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/test_type"
	workOrderRepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/work_order"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/usecase"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/result"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/util"
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

// formatNumberWithThousandSeparator formats a number with dot as thousand separator
func formatNumberWithThousandSeparator(value float64, decimals int) string {
	// For values that might cause issues, return simple format
	if value == 0 {
		return "0"
	}

	// Format the number with specified decimals
	formatted := fmt.Sprintf(fmt.Sprintf("%%.%df", decimals), value)

	// For very large numbers or edge cases, return the simple format
	if len(formatted) > 15 {
		return fmt.Sprintf(fmt.Sprintf("%%.%df", decimals), value)
	}

	// Split into integer and decimal parts
	parts := strings.Split(formatted, ".")
	integerPart := parts[0]

	// Add thousand separators to integer part
	if len(integerPart) > 3 {
		var result strings.Builder
		for i, digit := range integerPart {
			if i > 0 && (len(integerPart)-i)%3 == 0 {
				result.WriteString(".")
			}
			result.WriteRune(digit)
		}
		integerPart = result.String()
	}

	// Reconstruct the number
	if decimals > 0 && len(parts) > 1 {
		return integerPart + "," + parts[1] // Use comma for decimal separator
	}
	return integerPart
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

		alias := testResult.TestType.AliasCode
		if alias == "" && len(testResult.History) > 0 {
			alias = testResult.History[0].TestType.AliasCode
		}

		barcode := workOrder.BarcodeSIMRS // ono for new, visit no for old
		if barcode == "" {
			continue
		}

		// Translate visit no to ono
		order, err := u.repo.GetLisOrderByOnoOrVisitNo(barcode, barcode)
		if err != nil {
			return fmt.Errorf("error getting lis order by visit no: %w", err)
		}
		barcode = order.ONO

		code := alias
		if code == "" {
			code = testResult.TestType.Code
		}

		// Convert result and unit for specific tests
		resultValue := testResult.GetResult()
		unit := testResult.Unit
		refRange := testResult.ReferenceRange

		// Parse the string result to float64 for calculations
		resultFloat, err := strconv.ParseFloat(resultValue, 64)
		if err != nil {
			slog.Error("failed to parse result value", "error", err, "resultValue", resultValue)
			resultFloat = 0 // Default to 0 if parsing fails
		}

		if strings.TrimSpace(alias) == "Jumlah Trombosit" || strings.TrimSpace(alias) == "Jumlah Leukosit" {
			resultFloat = resultFloat * 1000
			unit = "/µL"

			// Convert reference range by multiplying by 1000
			if testResult.ReferenceRange != "" {
				// Use decimal from TestType for proper formatting
				decimal := testResult.TestType.Decimal
				if decimal < 0 {
					decimal = 0
				}
				refRangeConverted := util.ConvertReferenceRangeWithDecimal(testResult.ReferenceRange, 1000, decimal)
				if refRangeConverted != "" {
					refRange = refRangeConverted
				}
			}
		}

		// Round the result value (e.g., 8.4 -> 8, 8.6 -> 9)
		// Exception: Don't round Hemoglobin and Eritrosit
		var resultValueStr string
		// Special formatting for Trombosit and Leukosit (already converted by 1000 above)
		if strings.TrimSpace(alias) == "Jumlah Trombosit" || strings.TrimSpace(alias) == "Jumlah Leukosit" {
			resultValueStr = formatNumberWithThousandSeparator(resultFloat, 0)
		} else {
			// Keep decimal for other tests with thousand separator
			resultValueStr = formatNumberWithThousandSeparator(resultFloat, 1)
		}
		// } else {
		// 	// Round other tests to whole numbers with thousand separator
		// 	roundedValue := math.Round(resultValue)
		// resultValueStr = formatNumberWithThousandSeparator(roundedValue, 0)
		// }

		// Log the values for debugging
		flag := entity.NewKhanzaFlag(testResult)

		// Ensure FLAG field is compatible with database constraints
		// Convert flag to single character or empty string
		var flagStr string
		switch flag {
		case entity.FlagHigh:
			flagStr = "H"
		case entity.FlagLow:
			flagStr = "L"
		case entity.FlagLowLow:
			flagStr = "L" // Use single L for compatibility
		case entity.FlagHighHigh:
			flagStr = "H" // Use single H for compatibility
		default:
			flagStr = "" // Empty for normal or no data
		}

		slog.Info("khanza sync debug",
			"test_code", code,
			"alias_code", alias,
			"original_value", resultValue,
			"formatted_value", resultValueStr,
			"formatted_length", len(resultValueStr),
			"flag", flag,
			"flag_str", flagStr,
			"flag_length", len(flagStr),
		)

		// Validate field lengths to prevent database truncation errors
		if len(resultValueStr) > 50 {
			resultValueStr = resultValueStr[:50] // Truncate if too long
		}
		if len(code) > 20 {
			code = code[:20] // Truncate test code if too long
		}
		if len(unit) > 10 {
			unit = unit[:10] // Truncate unit if too long
		}
		if len(refRange) > 100 {
			refRange = refRange[:100] // Truncate reference range if too long
		}

		khanzaDT := entity.KhanzaResDT{
			ONO:         barcode,
			TESTCD:      code,
			TestNM:      code,
			DataTyp:     entity.DataTypNumeric,
			ResultValue: resultValueStr,
			Unit:        unit,
			Flag:        entity.Flag(flagStr), // Use converted flag string
			RefRange:    refRange,
		}
		reqs = append(reqs, khanzaDT)
	}

	err = u.repo.BatchUpsertRESDTO(ctx, reqs)
	if err != nil {
		slog.Error("batch upsert failed",
			"error", err,
			"total_records", len(reqs),
			"work_order_id", workOrderID,
		)

		// Log sample of records for debugging
		for i, req := range reqs {
			if i < 3 { // Log first 3 records for debugging
				slog.Error("sample record",
					"index", i,
					"ono", req.ONO,
					"test_cd", req.TESTCD,
					"result_value", req.ResultValue,
					"result_value_len", len(req.ResultValue),
					"flag", req.Flag,
					"flag_len", len(string(req.Flag)),
					"unit", req.Unit,
					"unit_len", len(req.Unit),
				)
			}
		}

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
		BarcodeSIMRS: ono,
		TestTypes:    tests,
		CreatedBy:    -1,
	}

	existingWorkOrder, err := u.workOrderRepository.GetBySIMRSBarcode(ctx, ono)
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
	rawRequest = u.fixBrokenPName(rawRequest)
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

func (u *Usecase) GetResult(ctx context.Context, ono string) (Response, error) {
	slog.InfoContext(ctx, "debug khanza get result", "ono", ono)

	orders, err := u.repo.GetLabRequestByNoOrder(ctx, ono)
	if err != nil {
		return Response{}, fmt.Errorf("error getting order by no order: %w", err)
	}

	if len(orders) == 0 {
		return Response{}, fmt.Errorf("order not found")
	}

	firstOrders := orders[0]
	workOrder, err := u.getWorkOrderWithFallback(ctx, firstOrders.NoOrder, firstOrders.NoRawat)
	if err != nil {
		return Response{}, err
	}

	workOrder.FillTestResultDetail(false)
	resultTestMap := u.groupTestResultByAliasCode(workOrder.TestResult)
	resultTest := make([]ResponseResultTest, len(orders))

	for i, o := range orders {
		t, ok := resultTestMap[o.Pemeriksaan]
		if !ok {
			resultTest[i] = ResponseResultTest{
				TestID:      o.IDTemplate,
				NamaTest:    o.Pemeriksaan,
				Hasil:       "",
				NilaiNormal: o.NilaiRujukanPA,
				Satuan:      o.Satuan,
				Flag:        "",
			}
			continue
		}

		hasil := ""
		if t.Result != "" {
			// Parse the string result to float64, then format it
			if resultFloat, err := strconv.ParseFloat(t.Result, 64); err == nil {
				hasil = strconv.FormatFloat(resultFloat, 'f', t.TestType.Decimal, 64)
			} else {
				// If parsing fails, use the string value as is
				hasil = t.Result
			}
		}

		resultTest[i] = u.resultConvert(ResponseResultTest{
			TestID:      o.IDTemplate,
			NamaTest:    o.Pemeriksaan,
			Hasil:       hasil,
			NilaiNormal: t.ReferenceRange,
			Satuan:      t.Unit,
			Flag:        string(entity.NewKhanzaFlag(t)),
		})
	}

	res := Response{}
	res.Result.OBX.OrderLab = ono
	res.Response.Sample.ResultTest = resultTest
	slog.InfoContext(ctx, "debug khanza get result", "res", res)

	return res, nil
}

func (u *Usecase) getWorkOrderWithFallback(ctx context.Context, ono string, visitNo string) (entity.WorkOrder, error) {
	workOrder, err := u.workOrderRepository.GetBySIMRSBarcode(ctx, ono)
	if err != nil && !errors.Is(err, entity.ErrNotFound) {
		return entity.WorkOrder{}, fmt.Errorf("error getting work order by ono: %w", err)
	}

	if workOrder.ID != 0 {
		return workOrder, nil
	}

	workOrder, err = u.workOrderRepository.GetBySIMRSBarcode(ctx, visitNo)
	if err != nil {
		return entity.WorkOrder{}, fmt.Errorf("error getting work order by visit no: %w", err)
	}

	return workOrder, nil
}

func (*Usecase) resultConvert(result ResponseResultTest) ResponseResultTest {
	if strings.TrimSpace(result.NamaTest) == "Jumlah Trombosit" || strings.TrimSpace(result.NamaTest) == "Jumlah Leukosit" {
		// For GetResult method, we need to convert the display format
		// The values here are already in base units (17.8), so convert for display
		if result.Hasil != "" {
			// Parse the number value
			value, err := strconv.ParseFloat(result.Hasil, 64)
			if err == nil {
				// Multiply by 1000 and format with dot as thousand separator
				convertedValue := value * 1000
				result.Hasil = formatNumberWithThousandSeparator(convertedValue, 0)
			}
		}

		// Convert reference range values
		if strings.Contains(result.NilaiNormal, "-") {
			parts := strings.Split(result.NilaiNormal, "-")
			if len(parts) == 2 {
				minStr := strings.TrimSpace(parts[0])
				maxStr := strings.TrimSpace(parts[1])

				minVal, minErr := strconv.ParseFloat(minStr, 64)
				maxVal, maxErr := strconv.ParseFloat(maxStr, 64)

				if minErr == nil && maxErr == nil {
					minConverted := formatNumberWithThousandSeparator(minVal*1000, 0)
					maxConverted := formatNumberWithThousandSeparator(maxVal*1000, 0)
					result.NilaiNormal = fmt.Sprintf("%s - %s", minConverted, maxConverted)
				}
			}
		}

		result.Satuan = "uL"
	}

	if strings.TrimSpace(result.NamaTest) == "Glukosa Sewaktu" {
		result.NilaiNormal = "< 180"
	}

	if strings.TrimSpace(result.NamaTest) == "HDL - Kolesterol" {
		h, _ := strconv.Atoi(result.Hasil)
		if h < 40 {
			result.Flag = "H"
		} else if h >= 40 {
			result.Flag = ""
		}
		result.NilaiNormal = ">40"
	}

	if strings.TrimSpace(result.NamaTest) == "LDL - Kolesterol" {
		result.NilaiNormal = "<100"
	}

	if strings.TrimSpace(result.NamaTest) == "Total - Kolesterol" {
		result.NilaiNormal = "<200"
	}

	if strings.TrimSpace(result.NamaTest) == "Trigliserida" {
		result.NilaiNormal = "<150"
	}

	if strings.TrimSpace(result.NamaTest) == "Eritrosit" {
		result.Satuan = "10^6/uL"
	}

	if strings.TrimSpace(result.NamaTest) == "Bilirubin Direk" {
		result.NilaiNormal = "< 0,5"
	}

	if strings.TrimSpace(result.NamaTest) == "Asam Urat" {
		result.NilaiNormal = "< 5,7"
	}

	// Convert decimal separator from dot to comma for most tests, but keep dots for thousand separators in Trombosit/Leukosit
	if strings.TrimSpace(result.NamaTest) == "Jumlah Trombosit" || strings.TrimSpace(result.NamaTest) == "Jumlah Leukosit" {
		// Keep dots as thousand separators for these tests (250.000, not 250,000)
	} else {
		// Convert decimal separator from dot to comma for other tests (5.4 -> 5,4)
		result.Hasil = strings.ReplaceAll(result.Hasil, ".", ",")
	}

	// Remove unnecessary zeros from NilaiNormal (exclude specific tests)
	if strings.TrimSpace(result.NamaTest) == "Jumlah Trombosit" || strings.TrimSpace(result.NamaTest) == "Jumlah Leukosit" || strings.TrimSpace(result.NamaTest) == "Eritrosit" {
		// Keep NilaiNormal as is for these tests
	} else {
		// Remove unnecessary zeros from NilaiNormal throughout the string (2.00 - 4.00 -> 2 - 4)
		result.NilaiNormal = strings.ReplaceAll(result.NilaiNormal, ".00", "")
		result.NilaiNormal = strings.ReplaceAll(result.NilaiNormal, ".0", "")
		result.NilaiNormal = strings.ReplaceAll(result.NilaiNormal, ",00", "")
		result.NilaiNormal = strings.ReplaceAll(result.NilaiNormal, ",0", "")
	}

	// Remove unnecessary zeros from Hasil (exclude Trombosit/Leukosit to preserve thousand separators)
	if strings.TrimSpace(result.NamaTest) == "Jumlah Trombosit" || strings.TrimSpace(result.NamaTest) == "Jumlah Leukosit" {
		// Keep thousand separators for these tests (250.000, not 250)
	} else {
		// Remove unnecessary zeros for other tests (12.0 -> 12, but keep 12.05 as 12,05)
		result.Hasil = strings.TrimSuffix(result.Hasil, ".0")
		result.Hasil = strings.TrimSuffix(result.Hasil, ",0")
	}

	return result
}

func (u *Usecase) groupedByTestName(ctx context.Context, orders []entity.KhanzaLabRequest) (map[string]entity.KhanzaLabRequest, error) {
	grouped := make(map[string]entity.KhanzaLabRequest)
	for _, order := range orders {
		grouped[order.Pemeriksaan] = order
	}
	return grouped, nil
}

func (u *Usecase) groupTestResultByAliasCode(orders []entity.TestResult) map[string]entity.TestResult {
	grouped := make(map[string]entity.TestResult)
	for _, order := range orders {
		testName := order.TestType.AliasCode
		if testName == "" {
			testName = order.TestType.Name
		}
		grouped[testName] = order
	}
	return grouped
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

func (u *Usecase) fixBrokenPName(raw []byte) []byte {
	re := regexp.MustCompile(`"pname"\s*:\s*"([^"]*?)"([^",}]+)"`)
	return re.ReplaceAll(raw, []byte(`"pname": "$1'$2"`))
}
