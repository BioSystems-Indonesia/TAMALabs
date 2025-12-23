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

	obxParts := strings.Split(obx, "|")
	if len(obxParts) < 6 {
		return fmt.Errorf("invalid OBX segment")
	}

	testCode := strings.Split(obxParts[3], "^")[0]
	measuredValue, err := strconv.ParseFloat(obxParts[5], 64)
	if err != nil {
		return fmt.Errorf("invalid measured value")
	}

	operator := "SYSTEM"
	if len(obxParts) > 16 {
		operator = obxParts[16]
	}

	device, err := u.deviceRepo.FindOneByType(ctx, entity.DeviceType(deviceIdentifier))
	if err != nil {
		return err
	}

	testType, err := u.testTypeRepo.FindOneByCode(ctx, testCode)
	if err != nil {
		return err
	}

	qcLevel := u.extractQCLevel(spm)
	qcEntry, err := u.qcEntryRepo.GetActiveEntry(ctx, device.ID, testType.ID, qcLevel)
	if err != nil {
		return err
	}

	if qcEntry.TargetSD == nil || *qcEntry.TargetSD == 0 {
		return fmt.Errorf("target SD not defined for QC entry")
	}

	prevStats, _ := u.qcResultRepo.GetByEntryIDAndMethod(
		ctx,
		qcEntry.ID,
		"statistic",
	)

	values := make([]float64, 0, len(prevStats)+1)
	for _, r := range prevStats {
		values = append(values, r.MeasuredValue)
	}
	values = append(values, measuredValue)

	manualMean := qcEntry.TargetMean
	manualSD := *qcEntry.TargetSD
	manualCV := (manualSD / manualMean) * 100
	manualErrorSD := (measuredValue - manualMean) / manualSD

	manualResult := u.buildQCResult(
		qcEntry.ID,
		measuredValue,
		manualMean,
		manualSD,
		manualCV,
		manualErrorSD,
		"manual",
		operator,
		messageControlID,
	)

	if err := u.qcResultRepo.Create(ctx, manualResult); err != nil {
		return err
	}

	var statMean, statSD, statCV float64

	targetMean := manualMean
	targetSD := manualSD

	count := len(values)

	if count < 5 {
		statMean = targetMean
		statSD = targetSD
		statCV = manualCV

	} else {
		stats := calculateStats(values)

		statMean = stats.Mean
		statSD = stats.SD

		minSD := targetSD * 0.7
		if statSD < minSD {
			statSD = minSD
		}

		alpha := math.Min(float64(count-5)/5.0, 1.0)

		statMean = alpha*statMean + (1-alpha)*targetMean
		statSD = alpha*statSD + (1-alpha)*targetSD

		if statMean != 0 {
			statCV = (statSD / statMean) * 100
		}
	}

	statErrorSD := 0.0
	if statSD != 0 {
		statErrorSD = (measuredValue - statMean) / statSD
	}

	statResult := u.buildQCResult(
		qcEntry.ID,
		measuredValue,
		statMean,
		statSD,
		statCV,
		statErrorSD,
		"statistic",
		operator,
		messageControlID,
	)

	return u.qcResultRepo.Create(ctx, statResult)
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
	// Use QCEntry repository to gather device summary
	summary, err := u.qcEntryRepo.GetDeviceSummary(ctx, deviceID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_qc":         summary.TotalQC,
		"qc_this_month":    summary.QCThisMonth,
		"last_qc":          summary.LastQCDate,
		"last_qc_status":   summary.LastQCStatus,
		"qc_today_status":  summary.QCTodayStatus,
		"level_1_complete": summary.Level1Complete,
		"level_2_complete": summary.Level2Complete,
		"level_3_complete": summary.Level3Complete,
		"level_1_today":    summary.Level1Today,
		"level_2_today":    summary.Level2Today,
		"level_3_today":    summary.Level3Today,
	}, nil
}

type qcStats struct {
	Mean float64
	SD   float64
	CV   float64
}

type StatsResult struct {
	Mean float64
	SD   float64
	CV   float64
}

func calculateStats(values []float64) StatsResult {
	n := len(values)
	if n == 0 {
		return StatsResult{}
	}

	var sum float64
	for _, v := range values {
		sum += v
	}
	mean := sum / float64(n)

	var variance float64
	for _, v := range values {
		diff := v - mean
		variance += diff * diff
	}

	sd := 0.0
	if n > 1 {
		sd = math.Sqrt(variance / float64(n-1))
	}

	cv := 0.0
	if mean != 0 {
		cv = (sd / mean) * 100
	}

	return StatsResult{
		Mean: mean,
		SD:   sd,
		CV:   cv,
	}
}

func (u *QualityControlUsecase) buildQCResult(
	qcEntryID int,
	measured, mean, sd, cv, errorSD float64,
	method string,
	operator, messageControlID string,
) *entity.QCResult {

	absError := measured - mean
	relativeError := 0.0
	if mean != 0 {
		relativeError = (absError / mean) * 100
	}

	result := "In Control"
	if errorSD >= 3 || errorSD <= -3 {
		result = "Reject"
	} else if (errorSD > 2 && errorSD < 3) || (errorSD < -2 && errorSD > -3) {
		result = "Warning"
	}

	return &entity.QCResult{
		QCEntryID:     qcEntryID,
		MeasuredValue: measured,

		CalculatedMean: mean,
		CalculatedSD:   sd,
		CalculatedCV:   cv,
		ErrorSD:        errorSD,

		SD1:              mean + sd,
		SD2:              mean + 2*sd,
		SD3:              mean + 3*sd,
		AbsoluteError:    absError,
		RelativeError:    relativeError,
		Result:           result,
		Method:           method,
		Operator:         operator,
		MessageControlID: messageControlID,
	}
}
