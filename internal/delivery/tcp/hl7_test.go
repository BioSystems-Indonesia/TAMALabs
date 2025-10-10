package tcp

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/mock"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/usecase"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const ex = `MSH|^~\&|BA200|Biosystems|Host|Host provider|20241213131746||OUL^R22^OUL_R22|a5f12c06-2d66-4a2e-b1bd-f0423978db14|P|2.5.1|||ER|AL||UNICODE UTF-8|||LAB-29^IHE
SPM|1|SI APA||SER|||||||P|||||||
OBR||""||ACE^ACE^A400|||||||||||||||||||||||||
ORC|OK||||CM||||20241213131746
OBX|1|NM|ACE^ACE^A400||21.9271355|U/L^U/L^A400||NONE|||F|||||ADMIN||A400^Biosystems~834000237^Biosystems|20241213193435
OBR||""||ALBUMIN^ALBUMIN^A400|||||||||||||||||||||||||
ORC|OK||||CM||||20241213131746
OBX|1|NM|ALBUMIN^ALBUMIN^A400||43.3798141|g/L^g/L^A400||NONE|||F|||||ADMIN||A400^Biosystems~834000237^Biosystems|20241213192844
OBR||""||ALP-AMP^ALP-AMP^A400|||||||||||||||||||||||||
ORC|OK||||CM||||20241213131746
OBX|1|NM|ALP-AMP^ALP-AMP^A400||35.8951797|U/L^U/L^A400||NONE|||F|||||ADMIN||A400^Biosystems~834000237^Biosystems|20241213193520
OBR||""||ALT-GPT^ALT-GPT^A400|||||||||||||||||||||||||
ORC|OK||||CM||||20241213131746
OBX|1|NM|ALT-GPT^ALT-GPT^A400||10.5633411|U/L^U/L^A400||NONE|||F|||||ADMIN||A400^Biosystems~834000237^Biosystems|20241213193538
OBR||""||AMYLASE DIRECT^AMYLASE DIRECT^A400|||||||||||||||||||||||||
ORC|OK||||CM||||20241213131746
OBX|1|NM|AMYLASE DIRECT^AMYLASE DIRECT^A400||106.056778|U/L^U/L^A400||NONE|||F|||||ADMIN||A400^Biosystems~834000237^Biosystems|20241213201146
OBR||""||AST-GOT^AST-GOT^A400|||||||||||||||||||||||||
ORC|OK||||CM||||20241213131746
OBX|1|NM|AST-GOT^AST-GOT^A400||21.2964096|U/L^U/L^A400||NONE|||F|||||ADMIN||A400^Biosystems~834000237^Biosystems|20241213201728
OBR||""||BILI DIRECT DPD^BILI DIRECT DPD^A400|||||||||||||||||||||||||
ORC|OK||||CM||||20241213131746
OBX|1|NM|BILI DIRECT DPD^BILI DIRECT DPD^A400||0.0581453256|mg/dL^mg/dL^A400||NONE|||F|||||ADMIN||A400^Biosystems~834000237^Biosystems|20241213193238
OBR||""||BILI TOTAL DPD^BILI TOTAL DPD^A400|||||||||||||||||||||||||
ORC|OK||||CM||||20241213131746
OBX|1|NM|BILI TOTAL DPD^BILI TOTAL DPD^A400||0.71246928|mg/dL^mg/dL^A400||NONE|||F|||||ADMIN||A400^Biosystems~834000237^Biosystems|20241213193511
OBR||""||CALCIUM ARSENAZO^CALCIUM ARSENAZO^A400|||||||||||||||||||||||||
ORC|OK||||CM||||20241213131746
OBX|1|NM|CALCIUM ARSENAZO^CALCIUM ARSENAZO^A400||9.38908291|mg/dL^mg/dL^A400||NONE|||F|||||ADMIN||A400^Biosystems~834000237^Biosystems|20241213193229
OBR||""||CHOLESTEROL^CHOLESTEROL^A400|||||||||||||||||||||||||
ORC|OK||||CM||||20241213131746
OBX|1|NM|CHOLESTEROL^CHOLESTEROL^A400||176.641998|mg/dL^mg/dL^A400||NONE|||F|||||ADMIN||A400^Biosystems~834000237^Biosystems|20241213193220
OBR||""||CHOLINESTERASE^CHOLINESTERASE^A400|||||||||||||||||||||||||
ORC|OK||||CM||||20241213131746
OBX|1|NM|CHOLINESTERASE^CHOLINESTERASE^A400||6649.88037|U/L^U/L^A400||NONE|||F|||||ADMIN||A400^Biosystems~834000237^Biosystems|20241213193435
OBR||""||CK^CK^A400|||||||||||||||||||||||||
ORC|OK||||CM||||20241213131746
OBX|1|NM|CK^CK^A400||22.9780865|U/L^U/L^A400||NONE|||F|||||ADMIN||A400^Biosystems~834000237^Biosystems|20241213193453
OBR||""||CREATININE^CREATININE^A400|||||||||||||||||||||||||
ORC|OK||||CM||||20241213131746
OBX|1|NM|CREATININE^CREATININE^A400||0.635651588|mg/dL^mg/dL^A400||NONE|||F|||||ADMIN||A400^Biosystems~834000237^Biosystems|20241213193332
OBR||""||FERRITIN^FERRITIN^A400|||||||||||||||||||||||||
ORC|OK||||CM||||20241213131746
OBX|1|NM|FERRITIN^FERRITIN^A400||7.29504681|\Z03BC\g/L^\Z03BC\g/L^A400||NONE|||F|||||ADMIN||A400^Biosystems~834000237^Biosystems|20241213193520
OBR||""||GAMMA-GT^GAMMA-GT^A400|||||||||||||||||||||||||
ORC|OK||||CM||||20241213131746
OBX|1|NM|GAMMA-GT^GAMMA-GT^A400||14.6088457|U/L^U/L^A400||NONE|||F|||||ADMIN||A400^Biosystems~834000237^Biosystems|20241213193435
OBR||""||GLUCOSE^GLUCOSE^A400|||||||||||||||||||||||||
ORC|OK||||CM||||20241213131746
OBX|1|NM|GLUCOSE^GLUCOSE^A400||77.7815018|mg/dL^mg/dL^A400||NONE|||F|||||ADMIN||A400^Biosystems~834000237^Biosystems|20241213193511
OBR||""||HDL DIRECT TOOS^HDL DIRECT TOOS^A400|||||||||||||||||||||||||
ORC|OK||||CM||||20241213131746
OBX|1|NM|HDL DIRECT TOOS^HDL DIRECT TOOS^A400||79.2626419|mg/dL^mg/dL^A400||NONE|||F|||||ADMIN||A400^Biosystems~834000237^Biosystems|20241213193502
OBR||""||IRON FERROZINE^IRON FERROZINE^A400|||||||||||||||||||||||||
ORC|OK||||CM||||20241213131746
OBX|1|NM|IRON FERROZINE^IRON FERROZINE^A400||6.75784826|\Z03BC\g/dL^\Z03BC\g/dL^A400||NONE|||F|||||ADMIN||A400^Biosystems~834000237^Biosystems|20241213193511
OBR||""||LDL DIRECT TOOS^LDL DIRECT TOOS^A400|||||||||||||||||||||||||
ORC|OK||||CM||||20241213131746
OBX|1|NM|LDL DIRECT TOOS^LDL DIRECT TOOS^A400||95.1346283|mg/dL^mg/dL^A400||NONE|||F|||||ADMIN||A400^Biosystems~834000237^Biosystems|20241213193511
OBR||""||MAGNESIUM^MAGNESIUM^A400|||||||||||||||||||||||||
ORC|OK||||CM||||20241213131746
OBX|1|NM|MAGNESIUM^MAGNESIUM^A400||1.69774115|mg/dL^mg/dL^A400||NONE|||F|||||ADMIN||A400^Biosystems~834000237^Biosystems|20241213201746
OBR||""||PROTEIN TOTALBIR^PROTEIN TOTALBIR^A400|||||||||||||||||||||||||
ORC|OK||||CM||||20241213131746
OBX|1|NM|PROTEIN TOTALBIR^PROTEIN TOTALBIR^A400||70.7569351|g/L^g/L^A400||NONE|||F|||||ADMIN||A400^Biosystems~834000237^Biosystems|20241213193529
OBR||""||TRIGLYCERIDES^TRIGLYCERIDES^A400|||||||||||||||||||||||||
ORC|OK||||CM||||20241213131746
OBX|1|NM|TRIGLYCERIDES^TRIGLYCERIDES^A400||33.0238266|mg/dL^mg/dL^A400||NONE|||F|||||ADMIN||A400^Biosystems~834000237^Biosystems|20241213193220
OBR||""||UREA-BUN-UV^UREA-BUN-UV^A400|||||||||||||||||||||||||
ORC|OK||||CM||||20241213131746
OBX|1|NM|UREA-BUN-UV^UREA-BUN-UV^A400||24.036623|mg/dL^mg/dL^A400||NONE|||F|||||ADMIN||A400^Biosystems~834000237^Biosystems|20241213193323
OBR||""||URIC ACID^URIC ACID^A400|||||||||||||||||||||||||
ORC|OK||||CM||||20241213131746
OBX|1|NM|URIC ACID^URIC ACID^A400||2.88172531|mg/dL^mg/dL^A400||NONE|||F|||||ADMIN||A400^Biosystems~834000237^Biosystems|20241213193229`

const exQuery = `MSH|^~\\&|BA200|Biosystems|Host|Host provider|20250305201642||QBP^Q11^QBP_Q11|dd980fa0-6a3d-4584-8fee-5392f42eca1d|P|2.5.1|||ER|AL||UNICODE UTF-8|||LAB-27^IHE
QPD|WOS^Work Order Step^IHE_LABTF|dd980fa06a3d45848fee5392f42eca1d|20250107000003
RCP|I||R`

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
		// 	name: "QBP_Q11",
		// 	fields: fields{
		// 		AnalyzerUsecase: &analyzer.Usecase{
		// 			ObservationRequestRepository: &observation_request.Repository{
		// 				DB: db,
		// 			},
		// 			ObservationResultRepository: &observation_result.Repository{
		// 				DB: db,
		// 			},
		// 		},
		// 	},
		// 	args: args{
		// 		message: exQuery,
		// 	},
		// 	want: "QBP_Q11 processed",
		// },
		{
			name: "ORM_O01",
			fields: fields{
				AnalyzerUsecase: func(ctrl *gomock.Controller, t *testing.T) usecase.Analyzer {
					mock := mock.NewMockAnalyzer(ctrl)
					mock.EXPECT().ProcessORMO01(gomock.Any(), gomock.Any()).Return([]entity.Specimen{
						{
							ID:      1,
							Barcode: "WBL2506100001",
							Type:    "WBL",
							Patient: entity.Patient{
								ID:        1,
								FirstName: "John",
								LastName:  "Doe",
								Birthdate: time.Now().AddDate(-21, 0, 0),
								Sex:       "M",
							},
							ObservationRequest: []entity.ObservationRequest{
								{
									TestCode:        "WBC",
									TestDescription: "WBC",
									RequestedDate:   time.Now(),
									SpecimenID:      1,
								},
							},
						},
					}, nil)
					return mock
				},
			},
			args: args{
				message: replaceNewline(ormO01Query),
			},
			want: "ORM_O01 processed",
		},
		// {
		// 	name: "ORU_R01",
		// 	fields: fields{
		// 		AnalyzerUsecase: func(ctrl *gomock.Controller, t *testing.T) usecase.Analyzer {
		// 			mock := mock.NewMockAnalyzer(ctrl)
		// 			mock.EXPECT().ProcessORUR01(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, o entity.ORU_R01) error {
		// 				expected := entity.ORU_R01{}
		// 				// ignore MSH
		// 				expected.MSH = o.MSH
		// 				assert.Len(t, o.Patient, 1)
		// 				assert.NotEmpty(t, o.Patient[0].LastName)
		// 				assert.Len(t, o.Patient[0].Specimen, 1)
		// 				assert.NotEmpty(t, o.Patient[0].Specimen[0].Barcode)
		// 				assert.Len(t, o.Patient[0].Specimen[0].ObservationResult, 42)
		// 				for _, r := range o.Patient[0].Specimen[0].ObservationResult {
		// 					assert.NotEmpty(t, r.TestCode, "TestCode is empty for specimen %s", o.Patient[0].Specimen[0].Barcode)
		// 				}

		// 				return nil
		// 			})
		// 			return mock
		// 		},
		// 	},
		// 	args: args{
		// 		message: replaceNewline(orur01Query),
		// 	},
		// 	want: "MSH|^~\\&|||LIS|Lab01|20250615215511||ACK^R01^ACK|391|P|2.5.1|||ER|AL|ID|UNICODE UTF-8|||LAB-28^IHE\rMSA|AA|391|Message accepted|||",
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			h := &HlSevenHandler{
				analyzerUsecase: tt.fields.AnalyzerUsecase(ctrl, t),
			}
			got, err := h.HL7Handler(context.TODO(), tt.args.message)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func replaceNewline(ormO01Query string) string {
	return strings.ReplaceAll(ormO01Query, "\n", "\r")
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
