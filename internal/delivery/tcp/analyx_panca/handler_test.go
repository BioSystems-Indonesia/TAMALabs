package analyxpanca

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"testing"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/mock"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/usecase"
	"github.com/BioSystems-Indonesia/TAMALabs/pkg/mllp"
	"github.com/BioSystems-Indonesia/TAMALabs/pkg/server"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const ormO01Query = "MSH|^~\\&|||||20250610192807||ORM^O01|2|P|2.3.1||||||UTF-8|||\rORC|RF|C1|wbl2506100001||IP\rOBX|1|IS|^BothLis Switch^||O||||||F|||||||\rOBX|2|IS|^Take Mode^||||||||F|||||||"
const orur01Query = "MSH|^~\\&|||||20250610192645||ORU^R01|391|P|2.3.1|||||CHA|UTF-8|||\rPID|1||9||dr.Barir|||O|||||||||||||||||||||||\rPV1|1|||||||||||||||||||||\rOBR|1||testbarcode|||20250610152043|20250610152328|||||||20250610145043||||Admin||||||HM||||||||||||||||\rOBX|1|NM|WBC||7.56|10^9/L|4.00-10.00||||F|||||||\rOBX|2|NM|Neu#||4.34|10^9/L|2.00-7.00||||F|||||||\rOBX|3|NM|Lym#||2.39|10^9/L|0.80-4.00||||F|||||||\rOBX|4|NM|Mon#||0.45|10^9/L|0.12-1.20||||F|||||||\rOBX|5|NM|Eos#||0.36|10^9/L|0.02-0.50||||F|||||||\rOBX|6|NM|Bas#||0.02|10^9/L|0.00-0.10||||F|||||||\rOBX|7|NM|Neu%||57.3|%|50.0-70.0||||F|||||||\rOBX|8|NM|Lym%||31.6|%|20.0-40.0||||F|||||||\rOBX|9|NM|Mon%||6.0|%|3.0-12.0||||F|||||||\rOBX|10|NM|Eos%||4.8|%|0.5-5.0||||F|||||||\rOBX|11|NM|Bas%||0.3|%|0.0-1.0||||F|||||||\rOBX|12|NM|NLR||1.82||-||||F|||||||\rOBX|13|NM|PLR||178.66||-||||F|||||||\rOBX|14|NM|RBC||4.94|10^12/L|3.50-5.50||||F|||||||\rOBX|15|NM|HGB||13.6|g/dL|11.0-16.0||||F|||||||\rOBX|16|NM|HCT||41.8|%|37.0-54.0||||F|||||||\rOBX|17|NM|MCV||84.6|fL|80.0-100.0||||F|||||||\rOBX|18|NM|MCH||27.6|pg|27.0-34.0||||F|||||||\rOBX|19|NM|MCHC||326|g/L|320-360||||F|||||||\rOBX|20|NM|RDW-CV||12.7|%|11.0-16.0||||F|||||||\rOBX|21|NM|RDW-SD||38.9|fL|35.0-56.0||||F|||||||\rOBX|22|NM|PLT||427|10^9/L|100-300|H|||F|||||||\rOBX|23|NM|MPV||7.6|fL|7.0-11.0||||F|||||||\rOBX|24|NM|PDW-CV||12.7|%|9.0-17.0||||F|||||||\rOBX|25|NM|PDW-SD||10.4|fL|9.0-17.0||||F|||||||\rOBX|26|NM|PCT||0.326|%|0.108-0.282|H|||F|||||||\rOBX|27|NM|P-LCC||59|10^9/L|30-90||||F|||||||\rOBX|28|NM|P-LCR||13.8|%|11.0-45.0||||F|||||||\rOBX|29|IS|Take Mode||O||||||F|||||||\rOBX|30|IS|Blood Mode||WH||||||F|||||||\rOBX|31|IS|Test Mode||CBC+DIFF||||||F|||||||\rOBX|32|IS|Low Mode||L-WBC/PLT||||||F|||||||\rOBX|33|IS|Ref Group||General||||||F|||||||\rOBX|34|IS|Age||||||||F|||||||\rOBX|35|IS|Remarks||||||||F|||||||\rOBX|36|IS|Blood Type||||||||F|||||||\rOBX|37|IS|ESR||||||||F|||||||\rOBX|38|IS|Recheck flag||Y||||||F|||||||\rOBX|39|IS|WBC Alarm||||||||F|||||||\rOBX|40|IS|RBC Alarm||||||||F|||||||\rOBX|41|IS|PLT Alarm||||||||F|||||||\rOBX|42|IS|Print_BMP||C||||||F|||||||\r"

func TestHlSevenHandler(t *testing.T) {
	type fields struct {
		AnalyzerUsecase func(ctrl *gomock.Controller, t *testing.T) usecase.Analyzer
	}
	type args struct {
		message string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// {
		// 	name: "ORM_O01",
		// 	fields: fields{
		// 		AnalyzerUsecase: func(ctrl *gomock.Controller, t *testing.T) usecase.Analyzer {
		// 			mock := mock.NewMockAnalyzer(ctrl)
		// 			mock.EXPECT().ProcessORMO01(gomock.Any(), gomock.Any()).Return([]entity.Specimen{
		// 				{
		// 					ID:      1,
		// 					Barcode: "WBL2506100001",
		// 					Type:    "WBL",
		// 					Patient: entity.Patient{
		// 						ID:        1,
		// 						FirstName: "John",
		// 						LastName:  "Doe",
		// 						Birthdate: time.Now().AddDate(-21, 0, 0),
		// 						Sex:       "M",
		// 					},
		// 					ObservationRequest: []entity.ObservationRequest{
		// 						{
		// 							TestCode:        "WBC",
		// 							TestDescription: "WBC",
		// 							RequestedDate:   time.Now(),
		// 							SpecimenID:      1,
		// 						},
		// 					},
		// 				},
		// 			}, nil)
		// 			return mock
		// 		},
		// 	},
		// 	args: args{
		// 		message: replaceNewline(ormO01Query),
		// 	},
		// 	want: "ORM_O01 processed",
		// },
		{
			name: "ORU_R01",
			fields: fields{
				AnalyzerUsecase: func(ctrl *gomock.Controller, t *testing.T) usecase.Analyzer {
					mock := mock.NewMockAnalyzer(ctrl)
					mock.EXPECT().ProcessORUR01(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, o entity.ORU_R01) error {
						t.Log("ORU_R01 decode")

						assert.Equal(t, "9", o.Patient[0].Specimen[0].Barcode)

						assert.Len(t, o.Patient, 1)
						assert.NotEmpty(t, o.Patient[0].LastName)
						assert.Len(t, o.Patient[0].Specimen, 1)
						assert.NotEmpty(t, o.Patient[0].Specimen[0].Barcode)
						assert.Len(t, o.Patient[0].Specimen[0].ObservationResult, 42)
						for _, r := range o.Patient[0].Specimen[0].ObservationResult {
							assert.NotEmpty(t, r.TestCode, "TestCode is empty for specimen %s", o.Patient[0].Specimen[0].Barcode)
						}

						return nil
					})
					return mock
				},
			},
			args: args{
				message: replaceNewline(orur01Query),
			},
			want: "MSH|^~\\&|||LIS|Lab01|20250615215511||ACK^R01^ACK|391|P|2.5.1|||ER|AL|ID|UNICODE UTF-8|||LAB-28^IHE\rMSA|AA|391|Message accepted|||",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			h := &Handler{
				analyzerUsecase: tt.fields.AnalyzerUsecase(ctrl, t),
			}

			// create a tcp server
			server := server.NewTCP("10523")
			server.SetHandler(h)
			err := server.Start()
			if err != nil {
				t.Fatal(err)
			}

			defer server.Stop()
			go server.Serve()

			conn, err := net.Dial("tcp", "127.0.0.1:10523")
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()

			mc := mllp.NewClient(conn)
			err = mc.Write([]byte(tt.args.message))
			if err != nil {
				t.Fatal(err)
			}

			res, err := mc.ReadAllRaw()
			if err != nil {
				t.Fatal(err)
			}

			// Write to response.bin
			err = os.WriteFile("response.bin", res, 0644)
			if err != nil {
				t.Fatal(err)
			}

			t.Logf("got: %v", replaceR(string(res)))
			// assert.Equal(t, tt.want, got)
		})
	}
}

func replaceNewline(ormO01Query string) string {
	return strings.ReplaceAll(ormO01Query, "\n", "\r")
}

func replaceR(message string) string {
	return strings.ReplaceAll(message, "\r", "\n")
}

func initSQLiteDB() (*gorm.DB, error) {

	// Change the directory using a relative path
	dbName := "../../../tmp/biosystem-lims.db"

	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	db.Logger = db.Logger.LogMode(logger.Silent)

	return db, err
}
