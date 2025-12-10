package test_template

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/config"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/util"
	"gorm.io/gorm"
)

type Repository struct {
	DB  *gorm.DB
	cfg *config.Schema
}

// NewRepository creates a new test type repository.
func NewRepository(db *gorm.DB, cfg *config.Schema) *Repository {
	return &Repository{DB: db, cfg: cfg}
}

// FindAll returns all test types.
func (r *Repository) FindAll(
	_ context.Context,
	req *entity.TestTemplateGetManyRequest,
) (entity.PaginationResponse[entity.TestTemplate], error) {
	db := r.DB.Preload("CreatedByUser").Preload("LastUpdatedByUser")
	sql.ProcessGetMany(db, req.GetManyRequest, sql.Modify{
		ProcessSearch: func(db *gorm.DB, query string) *gorm.DB {
			return db.Where("name like ?", "%"+query+"%").
				Or("description like ?", "%"+query+"%")
		}})

	resp, err := sql.GetWithPaginationResponse[entity.TestTemplate](db, req.GetManyRequest)
	if err != nil {
		return entity.PaginationResponse[entity.TestTemplate]{}, err
	}

	return resp, nil
}

func (r *Repository) FindOneByID(ctx context.Context, id int) (entity.TestTemplate, error) {
	db := r.DB.Preload("CreatedByUser").Preload("LastUpdatedByUser")
	var data entity.TestTemplate
	if err := db.First(&data, id).Error; err != nil {
		return entity.TestTemplate{}, err
	}

	return data, nil
}

func (r *Repository) Create(ctx context.Context, req *entity.TestTemplate) (entity.TestTemplate, error) {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(req).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return entity.TestTemplate{}, err
	}

	return *req, nil
}

func (r *Repository) GetObservationRequestDifference(ctx context.Context, req *entity.TestTemplate) (entity.TestTemplateObservationRequestDifference, error) {
	toDelete, toCreate, err := r.getToDeleteAndCreate(ctx, req)
	if err != nil {
		return entity.TestTemplateObservationRequestDifference{}, fmt.Errorf("error getting to delete and create: %w", err)
	}
	testTypeMapCode, err := r.getAllTestTypeCodeMapped(ctx)
	if err != nil {
		return entity.TestTemplateObservationRequestDifference{}, fmt.Errorf("error getting all test type code mapped: %w", err)
	}

	var workOrderTestTemplates []entity.WorkOrderTestTemplate
	if err := r.DB.Where("test_template_id = ?", req.ID).Find(&workOrderTestTemplates).Error; err != nil {
		return entity.TestTemplateObservationRequestDifference{}, fmt.Errorf("error finding work order test templates: %w", err)
	}

	var allObservationRequests []entity.ObservationRequest
	var toCreateObservationRequests []entity.ObservationRequest
	var toDeleteObservationRequests []entity.ObservationRequest
	for _, wo := range workOrderTestTemplates {
		var specimens []entity.Specimen
		if err := r.DB.Where("order_id =?", wo.WorkOrderID).Preload("WorkOrder").
			Preload("WorkOrder.Patient").
			Preload("ObservationRequest").Find(&specimens).Error; err != nil {
			return entity.TestTemplateObservationRequestDifference{}, fmt.Errorf("error finding specimens: %w", err)
		}

		groupedWorkOrderObservationRequest := make(map[int][]entity.ObservationRequest)
		for _, specimen := range specimens {
			groupedWorkOrderObservationRequest[int(wo.WorkOrderID)] = append(
				groupedWorkOrderObservationRequest[int(wo.WorkOrderID)],
				specimen.ObservationRequest...,
			)
		}

		for _, specimen := range specimens {
			for _, tt := range toCreate {
				if specimen.Type != tt.SpecimenType {
					continue
				}

				testType, ok := testTypeMapCode[tt.TestTypeCode]
				if !ok {
					slog.WarnContext(ctx, "test type not found", slog.Attr{
						Key:   "test_type_code",
						Value: slog.StringValue(tt.TestTypeCode),
					})
					continue
				}

				orGroup := groupedWorkOrderObservationRequest[int(wo.WorkOrderID)]
				if slices.ContainsFunc(orGroup, func(v entity.ObservationRequest) bool {
					return v.TestCode == testType.Code
				}) {
					continue
				}

				testTypeID := testType.ID
				newObservationRequest := entity.ObservationRequest{
					TestCode:        testType.Code,
					TestTypeID:      &testTypeID,
					TestDescription: testType.Name,
					SpecimenID:      int64(specimen.ID),
					RequestedDate:   time.Now(),
					WorkOrder:       specimen.WorkOrder,
				}
				toCreateObservationRequests = append(toCreateObservationRequests, newObservationRequest)
			}
		}

		for _, specimen := range specimens {
			// Backfill Specimen
			for j, _ := range specimen.ObservationRequest {
				specimen.ObservationRequest[j].WorkOrder = specimen.WorkOrder
			}

			allObservationRequests = append(allObservationRequests, specimen.ObservationRequest...)
		}
	}

	for _, testType := range toDelete {
		for _, observationRequest := range allObservationRequests {
			if observationRequest.TestCode == testType.TestTypeCode {

				toDeleteObservationRequests = append(toDeleteObservationRequests, observationRequest)
			}
		}
	}

	return entity.TestTemplateObservationRequestDifference{
		ToDelete: toDeleteObservationRequests,
		ToCreate: toCreateObservationRequests,
	}, nil
}

func (r *Repository) Update(ctx context.Context, req *entity.TestTemplate, diff *entity.TestTemplateObservationRequestDifference) (entity.TestTemplate, error) {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(req).Error; err != nil {
			return fmt.Errorf("error saving test template: %w", err)
		}

		if diff != nil {
			for _, observationRequest := range diff.ToDelete {
				if err := tx.Delete(&observationRequest).Error; err != nil {
					return fmt.Errorf("error deleting observation request: %w", err)
				}
			}

			for _, observationRequest := range diff.ToCreate {
				if err := tx.Create(&observationRequest).Error; err != nil {
					return fmt.Errorf("error creating observation request: %w", err)
				}
			}
		}

		return nil
	})
	if err != nil {
		return entity.TestTemplate{}, err
	}

	return *req, nil
}

func (r *Repository) getToDeleteAndCreate(ctx context.Context, req *entity.TestTemplate) ([]entity.WorkOrderCreateRequestTestType, []entity.WorkOrderCreateRequestTestType, error) {
	old, err := r.FindOneByID(ctx, req.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("test template not found: %v", err)
	}
	oldTestType := util.Map(old.RequestTestTypes, func(v entity.WorkOrderCreateRequestTestType) string {
		return fmt.Sprintf("%s:%s", v.TestTypeCode, v.SpecimenType)
	})
	newTestType := util.Map(req.RequestTestTypes, func(v entity.WorkOrderCreateRequestTestType) string {
		return fmt.Sprintf("%s:%s", v.TestTypeCode, v.SpecimenType)
	})

	toDelete, toCreate := util.CompareSlices(oldTestType, newTestType)
	toDeleteResp := util.Map(toDelete, func(v string) entity.WorkOrderCreateRequestTestType {
		parts := strings.Split(v, ":")
		return entity.WorkOrderCreateRequestTestType{
			TestTypeCode: parts[0],
			SpecimenType: parts[1],
		}
	})
	toCreateResp := util.Map(toCreate, func(v string) entity.WorkOrderCreateRequestTestType {
		parts := strings.Split(v, ":")
		return entity.WorkOrderCreateRequestTestType{
			TestTypeCode: parts[0],
			SpecimenType: parts[1],
		}
	})
	return toDeleteResp, toCreateResp, nil
}

func (r *Repository) getAllTestTypeCodeMapped(ctx context.Context) (map[string]entity.TestType, error) {
	var testTypes []entity.TestType
	if err := r.DB.WithContext(ctx).Find(&testTypes).Error; err != nil {
		return nil, err
	}
	testTypeMap := make(map[string]entity.TestType)
	for _, testType := range testTypes {
		testTypeMap[testType.Code] = testType
	}

	return testTypeMap, nil
}

func (r *Repository) Delete(ctx context.Context, req *entity.TestTemplate) (entity.TestTemplate, error) {
	if err := r.DB.Delete(req).Error; err != nil {
		return entity.TestTemplate{}, err
	}
	return *req, nil
}
