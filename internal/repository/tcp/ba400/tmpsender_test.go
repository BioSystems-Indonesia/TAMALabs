package ba400

import (
	"context"
	"testing"
	"time"

	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

func TestSendToBA400(t *testing.T) {
	t.Skip("need device")
	err := SendToBA400(context.Background(), []entity.Patient{
		{
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
							TestCode:        "ACID GLYCO BIR",
							TestDescription: "ACID GLYCO BIR",
							CreatedAt:       time.Now(),
							UpdatedAt:       time.Now(),
						},
						{
							ID:              2,
							SpecimenID:      1,
							TestCode:        "ALBUMIN",
							TestDescription: "ALBUMIN",
							CreatedAt:       time.Now(),
							UpdatedAt:       time.Now(),
						},
					},
				},
			},
		},
		{
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
							TestCode:        "ACID GLYCO BIR",
							TestDescription: "ACID GLYCO BIR",
							CreatedAt:       time.Now(),
							UpdatedAt:       time.Now(),
						},
						{
							ID:              4,
							SpecimenID:      2,
							TestCode:        "ALBUMIN",
							TestDescription: "ALBUMIN",
							CreatedAt:       time.Now(),
							UpdatedAt:       time.Now(),
						},
						{
							ID:              5,
							SpecimenID:      2,
							TestCode:        "URIC ACID",
							TestDescription: "URIC ACID",
							CreatedAt:       time.Now(),
							UpdatedAt:       time.Now(),
						},
					},
				},
			},
		},
	}, entity.Device{
		ID:        1,
		Name:      "Device 1",
		IPAddress: "192.168.1.100",
		Port:      5678,
	})

	if err != nil {
		t.Error(err)
	}
}

func Test_receiveResponse(t *testing.T) {
	type args struct {
		resp []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ORL_O34",
			args: args{
				resp: []byte(`MSH|^~\\&|BA200|Biosystems|Host|Host provider|20241223163505||ORL^O34^ORL_O34|ec3b41a9-77e3-4fd3-a2ca-f8f760dbda47|P|2.5.1|||ER|NE||UNICODE UTF-8|||LAB-28^IHE
MSA|AA|939b894f-a10a-4b35-9f82-95de095cc0c4`),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := receiveResponse(tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("receiveResponse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
