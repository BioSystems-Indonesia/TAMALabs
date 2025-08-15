package externaluc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
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

type UseCase struct {
}

func (u *UseCase) ProcessRequest(ctx context.Context, rawRequest []byte) error {
	var request Request
	var err error

	err = json.Unmarshal(rawRequest, &request)
	if err != nil {
		return fmt.Errorf("external.ProcessRequest json.Unmarshal failed: %w\nbody: %s", err, string(rawRequest))
	}

	var patient entity.Patient

	patient, err = u.findOrCreatePatient(request)

	w := &entity.WorkOrderCreateRequest{
		PatientID: patient.ID,
		CreatedBy: -100,
	}

	for _, testName := range request.Order.OBR.OrderTest {
		// TODO make this not N+1
		testType, err := u.TestTypeRepo.FindOneByCode(ctx, testName)
		if err != nil {
			slog.Error("find test type failed", "testName", testName, "error", err)
		}

		specimenType := "SER"
		if len(testType.Type) > 0 {
			specimenType = testType.Type[0].Type
		}

		w.TestTypes = append(w.TestTypes, entity.WorkOrderCreateRequestTestType{
			TestTypeID:   int64(testType.ID),
			TestTypeCode: testType.Code,
			SpecimenType: specimenType,
		})
	}

	_, err = u.WorkOrderRepo.Create(w)

	return err
}

func (u *UseCase) GetResult(ctx context.Context, id string) ([]byte, error) {
	res := Response{}

	workOrder, err := u.WorkOrderRepo.FindOneByBarcode(ctx, id)
	if err != nil {
		return nil, err
	}
	workOrder.FillTestResultDetail(false)

	res.Result.OBX.OrderLab = workOrder.Barcode

	resultTest := make([]ResponseResultTest, len(workOrder.TestResult))

	for i, t := range workOrder.TestResult {
		hasil := ""
		if t.Result != nil {
			hasil = strconv.FormatFloat(*t.Result, 'f', 4, 64)
		}
		resultTest[i] = ResponseResultTest{
			TestID:      strconv.FormatInt(t.ID, 10),
			NamaTest:    t.Test,
			Hasil:       hasil,
			NilaiNormal: t.ReferenceRange,
			Satuan:      t.Unit,
			Flag:        "",
		}
	}

	return json.Marshal(res)
}

func (u *UseCase) findOrCreatePatient(r Request) (entity.Patient, error) {
	p, err := u.convertIntoPatient(r)
	if err != nil {
		return p, err
	}

	return u.PatientRepo.FirstOrCreate(&p)
}

func (u *UseCase) convertIntoPatient(r Request) (entity.Patient, error) {

	firstName := r.Order.PID.PName
	lastName := ""

	names := strings.SplitN(r.Order.PID.PName, " ", 2)
	if len(names) == 2 {
		firstName = names[0]
		lastName = names[1]
	}

	sex := entity.PatientSexUnknown
	switch r.Order.PID.Sex {
	case "L":
		sex = entity.PatientSexMale
	case "F":
		sex = entity.PatientSexFemale
	}

	birthdate, err := time.Parse("2006-01-02", r.Order.PID.BirthDT)
	if err != nil {
		return entity.Patient{}, fmt.Errorf("cannot parsing birth date")
	}

	return entity.Patient{
		FirstName: firstName,
		LastName:  lastName,
		Sex:       sex,
		Birthdate: birthdate,
	}, nil
}

type Request struct {
	Order struct {
		PID struct {
			PName   string `json:"pname"`
			Sex     string `json:"sex"`
			BirthDT string `json:"birth_dt"`
		} `json:"pid"`
		OBR struct {
			OrderLab  string   `json:"order_lab"`
			OrderTest []string `json:"order_test"`
		}
	} `json:"order"`
}

type Response struct {
	Result struct {
		OBX struct {
			OrderLab string `json:"order_lab"`
		} `json:"obx"`
	} `json:"result"`

	Response struct {
		Sample struct {
			ResultTest []ResponseResultTest `json:"result_test"`
		} `json:"sampel"` // It's not typo its the request for
	} `json:"response"`
}

type ResponseResultTest struct {
	TestID      string `json:"test_id"`
	NamaTest    string `json:"nama_test"`
	Hasil       string `json:"hasil"`
	NilaiNormal string `json:"nilai_normal"`
	Satuan      string `json:"satuan"`
	Flag        string `json:"flag"`
}
