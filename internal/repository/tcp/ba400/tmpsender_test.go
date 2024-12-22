package ba400

import (
	"context"
	"testing"
	"time"

	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

func TestSendToBA400(t *testing.T) {
	err := SendToBA400(context.Background(), []entity.WorkOrderOMLRequest{
		{
			Patient: entity.Patient{
				ID:          1,
				FirstName:   "Pasien",
				LastName:    "Pertama",
				Birthdate:   time.Date(1995, time.January, 1, 0, 0, 0, 0, time.UTC),
				Sex:         "M",
				PhoneNumber: "",
				Location:    "",
				Address:     "",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Specimen: []entity.Specimen{
					{
						ID:             1,
						PatientID:      1,
						OrderID:        1,
						Type:           "SER",
						CollectionDate: time.Now().Format(time.RFC3339),
						ReceivedDate:   time.Time{},
						Source:         "",
						Condition:      "",
						Method:         "",
						Comments:       "",
						Barcode:        "123123",
						CreatedAt:      time.Now(),
						UpdatedAt:      time.Now(),
						ObservationRequest: []entity.ObservationRequest{
							{
								ID:              1,
								SpecimenID:      1,
								OrderID:         "1",
								TestCode:        "ACID GLYCO BIR",
								TestDescription: "ACID GLYCO BIR",
								CreatedAt:       time.Now(),
								UpdatedAt:       time.Now(),
							},
							{
								ID:              2,
								SpecimenID:      1,
								OrderID:         "1",
								TestCode:        "ALBUMIN",
								TestDescription: "ALBUMIN",
								CreatedAt:       time.Now(),
								UpdatedAt:       time.Now(),
							},
						},
					},
				},
			},
		},
		{
			Patient: entity.Patient{
				ID:          2,
				FirstName:   "Dafa",
				LastName:    "Kedua",
				Birthdate:   time.Date(2002, time.October, 23, 0, 0, 0, 0, time.UTC),
				Sex:         "F",
				PhoneNumber: "",
				Location:    "",
				Address:     "",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Specimen: []entity.Specimen{
					{
						ID:             2,
						PatientID:      2,
						OrderID:        1,
						Type:           "SER",
						CollectionDate: time.Now().Format(time.RFC3339),
						ReceivedDate:   time.Time{},
						Source:         "",
						Condition:      "",
						Method:         "",
						Comments:       "",
						Barcode:        "456456",
						CreatedAt:      time.Now(),
						UpdatedAt:      time.Now(),
						ObservationRequest: []entity.ObservationRequest{
							{
								ID:              3,
								SpecimenID:      2,
								OrderID:         "1",
								TestCode:        "ACID GLYCO BIR",
								TestDescription: "ACID GLYCO BIR",
								CreatedAt:       time.Now(),
								UpdatedAt:       time.Now(),
							},
							{
								ID:              4,
								SpecimenID:      2,
								OrderID:         "1",
								TestCode:        "ALBUMIN",
								TestDescription: "ALBUMIN",
								CreatedAt:       time.Now(),
								UpdatedAt:       time.Now(),
							},
							{
								ID:              5,
								SpecimenID:      2,
								OrderID:         "1",
								TestCode:        "URIC ACID",
								TestDescription: "URIC ACID",
								CreatedAt:       time.Now(),
								UpdatedAt:       time.Now(),
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		t.Error(err)
	}
}
