package quality_control

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"strconv"
	"strings"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/repository"
	devicerepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/device"
	test_type "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/test_type"
)

type QualityControlUsecase struct {
	qcEntryRepo  repository.QCEntry
	qcResultRepo repository.QCResult
	deviceRepo   *devicerepo.DeviceRepository
	testTypeRepo *test_type.Repository
}

func NewQualityControlUsecase(
	qcEntryRepo repository.QCEntry,
	qcResultRepo repository.QCResult,
	deviceRepo *devicerepo.DeviceRepository,
	testTypeRepo *test_type.Repository,
) *QualityControlUsecase {
	return &QualityControlUsecase{
		qcEntryRepo:  qcEntryRepo,
		qcResultRepo: qcResultRepo,
		deviceRepo:   deviceRepo,
		testTypeRepo: testTypeRepo,
	}
}

func (u *QualityControlUsecase) ParseAndSaveQC(ctx context.Context, hl7Message string, deviceIdentifier string) error {
	// Parse HL7 message
	lines := strings.Split(strings.ReplaceAll(hl7Message, "\r", "\n"), "\n")

	var msh, spm, inv, obx, orc string
	var messageControlID string

	for _, line := range lines {
		if strings.HasPrefix(line, "MSH|") {
			msh = line
			parts := strings.Split(line, "|")
			if len(parts) > 10 {
				messageControlID = parts[10]
			}
		} else if strings.HasPrefix(line, "SPM|") {
			spm = line
		} else if strings.HasPrefix(line, "INV|") {
			inv = line
		} else if strings.HasPrefix(line, "OBX|") {
			obx = line
		} else if strings.HasPrefix(line, "ORC|") {
			orc = line
		}

		if spm != "" && inv != "" && obx != "" {
			err := u.processQCRecord(ctx, msh, spm, inv, obx, orc, deviceIdentifier, messageControlID)
			if err != nil {
				slog.ErrorContext(ctx, "Error processing QC record", "error", err)
			}
			spm = ""
			inv = ""
			obx = ""
			orc = ""
		}
	}

	return nil
}

func (u *QualityControlUsecase) processQCRecord(
	ctx context.Context,
	msh, spm, inv, obx, orc, deviceIdentifier, messageControlID string,
) error {

	spmParts := strings.Split(spm, "|")
	if len(spmParts) < 3 {
		return fmt.Errorf("invalid SPM segment")
	}

	qcSampleID := spmParts[2]
	qcLevel := u.extractQCLevel(qcSampleID)

	_ = inv

	obxParts := strings.Split(obx, "|")
	if len(obxParts) < 8 {
		return fmt.Errorf("invalid OBX segment")
	}

	testCodeField := obxParts[3]
	testCodeParts := strings.Split(testCodeField, "^")
	if len(testCodeParts) < 1 {
		return fmt.Errorf("invalid test code")
	}
	testCode := testCodeParts[0]

	measuredMean, err := strconv.ParseFloat(obxParts[5], 64)
	if err != nil {
		return fmt.Errorf("invalid measured value: %w", err)
	}

	_ = obxParts[7]

	operator := "SYSTEM"
	if len(obxParts) > 16 {
		operator = obxParts[16]
	}

	orcParts := strings.Split(orc, "|")
	_ = orcParts // For future use

	deviceType := entity.DeviceType(deviceIdentifier)
	device, err := u.deviceRepo.FindOneByType(ctx, deviceType)
	if err != nil {
		return fmt.Errorf("device not found for type %s: %w", deviceIdentifier, err)
	}

	testType, err := u.testTypeRepo.FindOneByCode(ctx, testCode)
	if err != nil {
		return fmt.Errorf("test type not found for code %s: %w", testCode, err)
	}

	qcEntry, err := u.qcEntryRepo.GetActiveEntry(ctx, device.ID, testType.ID, qcLevel)
	if err != nil {
		return fmt.Errorf("no active QC entry found for device %d, test %d, level %d: %w",
			device.ID, testType.ID, qcLevel, err)
	}

	slog.InfoContext(ctx, "Debug", "measuredMean", measuredMean, "operator", operator, "qc entry", qcEntry)

	var targetMean, targetSD float64

	calculatedMean, _, count, err := u.qcResultRepo.CalculateStatistics(ctx, qcEntry.ID)
	if err != nil {
		slog.WarnContext(ctx, "Failed to calculate statistics, using target values from entry", "error", err)
		// Fallback to entry values
		targetMean = qcEntry.TargetMean
		if qcEntry.TargetSD != nil {
			targetSD = *qcEntry.TargetSD
		}
	} else if count >= 4 {

		tempMean := (calculatedMean*float64(count) + measuredMean) / float64(count+1)

		previousResults, _ := u.qcResultRepo.GetByEntryID(ctx, qcEntry.ID)
		var variance float64
		for _, r := range previousResults {
			diff := r.MeasuredValue - tempMean
			variance += diff * diff
		}
		diff := measuredMean - tempMean
		variance += diff * diff

		tempSD := 0.0
		totalCount := count + 1
		if totalCount > 1 {
			tempSD = math.Sqrt(variance / float64(totalCount-1))
		}

		targetMean = tempMean
		targetSD = tempSD
		slog.InfoContext(ctx, "Using calculated statistics (count >= 4, will be >= 5 after save)",
			"current_count", count,
			"count_after_save", count+1,
			"calculated_mean", tempMean,
			"calculated_sd", tempSD)
	} else {
		targetMean = qcEntry.TargetMean
		if qcEntry.TargetSD != nil {
			targetSD = *qcEntry.TargetSD
		}
		slog.InfoContext(ctx, "Using target values from entry (count < 4)",
			"current_count", count,
			"count_after_save", count+1)
	}

	var cvValue float64
	if targetMean != 0 {
		cvValue = (targetSD / targetMean) * 100
	}

	sd1 := targetMean + (1 * targetSD)
	sd2 := targetMean + (2 * targetSD)
	sd3 := targetMean + (3 * targetSD)

	var errorSD float64
	if qcEntry.TargetSD != nil && *qcEntry.TargetSD != 0 {
		errorSD = (measuredMean - qcEntry.TargetMean) / *qcEntry.TargetSD
	}

	absoluteError := measuredMean - qcEntry.TargetMean

	var relativeError float64
	if qcEntry.TargetMean != 0 {
		relativeError = (math.Abs(absoluteError) / qcEntry.TargetMean) * 100
	}

	var result string
	absErrorSD := math.Abs(errorSD)
	if absErrorSD <= 2 {
		result = "In Control"
	} else if absErrorSD <= 3 {
		result = "Warning"
	} else {
		result = "Reject"
	}

	// Create QC result record
	qcResult := entity.QCResult{
		QCEntryID:        qcEntry.ID,
		MeasuredValue:    measuredMean,
		CalculatedMean:   targetMean,
		CalculatedSD:     targetSD,
		CalculatedCV:     cvValue,
		ErrorSD:          errorSD,
		AbsoluteError:    absoluteError,
		RelativeError:    relativeError,
		SD1:              sd1,
		SD2:              sd2,
		SD3:              sd3,
		Result:           result,
		Method:           "statistic",
		Operator:         operator,
		MessageControlID: messageControlID,
	}

	// Save to database
	err = u.qcResultRepo.Create(ctx, &qcResult)
	if err != nil {
		return fmt.Errorf("error saving QC result: %w", err)
	}

	slog.InfoContext(ctx, "QC result saved",
		"device_id", device.ID,
		"test_type", testCode,
		"qc_level", qcLevel,
		"measured_value", measuredMean,
		"target_mean", targetMean,
		"target_sd", targetSD,
		"cv", cvValue,
		"result", result,
		"method", "statistic",
	)

	return nil
}

func (u *QualityControlUsecase) extractQCLevel(sampleID string) int {
	upper := strings.ToUpper(sampleID)
	if strings.Contains(upper, " III ") || strings.Contains(upper, " 3 ") {
		return 3
	}
	if strings.Contains(upper, " II ") || strings.Contains(upper, " 2 ") {
		return 2
	}
	if strings.Contains(upper, " I ") || strings.Contains(upper, " 1 ") {
		return 1
	}
	return 1 // Default to level 1
}

func (u *QualityControlUsecase) parseReferenceRange(refRange string) (ref, sd float64, err error) {
	parts := strings.Split(refRange, "-")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid reference range format")
	}

	min, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid min value: %w", err)
	}

	max, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid max value: %w", err)
	}

	ref = (min + max) / 2
	sd = (max - min) / 4 // Assuming range is ±2SD

	return ref, sd, nil
}

func (u *QualityControlUsecase) CreateQCEntry(ctx context.Context, req *entity.CreateQCEntryRequest) (*entity.QCEntry, error) {
	err := u.qcEntryRepo.DeactivateOldEntries(ctx, req.DeviceID, req.TestTypeID, req.QCLevel)
	if err != nil {
		return nil, fmt.Errorf("error deactivating old entries: %w", err)
	}

	entry := &entity.QCEntry{
		DeviceID:   req.DeviceID,
		TestTypeID: req.TestTypeID,
		QCLevel:    req.QCLevel,
		LotNumber:  req.LotNumber,
		TargetMean: req.TargetMean,
		TargetSD:   req.TargetSD,
		Method:     req.Method,
		IsActive:   true,
		RefMin:     req.RefMin,
		RefMax:     req.RefMax,
		CreatedBy:  req.CreatedBy,
	}

	err = u.qcEntryRepo.Create(ctx, entry)
	if err != nil {
		return nil, fmt.Errorf("error creating QC entry: %w", err)
	}

	return entry, nil
}

func (u *QualityControlUsecase) GetQCEntries(ctx context.Context, req entity.GetManyRequestQCEntry) ([]entity.QCEntry, int64, error) {
	return u.qcEntryRepo.GetMany(ctx, req)
}

func (u *QualityControlUsecase) GetQCResults(ctx context.Context, req entity.GetManyRequestQCResult) ([]entity.QCResult, int64, error) {
	return u.qcResultRepo.GetMany(ctx, req)
}

func (u *QualityControlUsecase) GetQCHistory(ctx context.Context, req entity.GetManyRequestQualityControl) ([]entity.QualityControl, int64, error) {
	// Return empty for backward compatibility
	return []entity.QualityControl{}, 0, nil
}

func (u *QualityControlUsecase) GetQCStatistics(ctx context.Context, deviceID int) (map[string]interface{}, error) {
	// Return empty for backward compatibility
	return map[string]interface{}{}, nil
}
