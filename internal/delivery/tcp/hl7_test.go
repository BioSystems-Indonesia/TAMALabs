package tcp

import (
	"context"
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/observation_request"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/observation_result"
	"github.com/oibacidem/lims-hl-seven/internal/usecase/analyzer"
	"github.com/stretchr/testify/assert"
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

func TestHlSevenHandler(t *testing.T) {
	db, _ := initSQLiteDB()
	type fields struct {
		AnalyzerUsecase *analyzer.Usecase
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
		{
			name: "QBP_Q11",
			fields: fields{
				AnalyzerUsecase: &analyzer.Usecase{
					ObservationRequestRepository: &observation_request.Repository{
						DB: db,
					},
					ObservationResultRepository: &observation_result.Repository{
						DB: db,
					},
				},
			},
			args: args{
				message: exQuery,
			},
			want: "QBP_Q11 processed",
		},
		{
			name: "OUL_R22_ELGATAMA",
			fields: fields{
				AnalyzerUsecase: &analyzer.Usecase{
					ObservationRequestRepository: &observation_request.Repository{
						DB: db,
					},
					ObservationResultRepository: &observation_result.Repository{
						DB: db,
					},
				},
			},
			args: args{
				message: "MSH|^~\\&|BA200|Biosystems|Host|Host provider|20241209214354||OUL^R22^OUL_R22|fac365d1-a158-4ce0-97f9-95afa1359571|P|2.5.1|||ER|AL||UNICODE UTF-8|||LAB-29^IH\rOBX|1|NM|UREA-BUN-UV^UREA-BUN-UV^A400||60.7323189|mg/dL^mg/dL^A400||NONE|||F|||||ADMIN||A400^Biosystems~834000237^Biosystems|20241210044024",
			},
			want: "OBX processed",
		},
		{
			name: "OUL_R22",
			fields: fields{
				AnalyzerUsecase: &analyzer.Usecase{
					ObservationRequestRepository: &observation_request.Repository{
						DB: db,
					},
					ObservationResultRepository: &observation_result.Repository{
						DB: db,
					},
				},
			},
			args: args{
				message: "MSH|^~\\&|LabSystem|LabFacility|EHR|Hospital|20241209120000||OUL^R22|MSG001|P|2.5.1\rPID|12|123|123456||John^Doe||19800101|M|||123 Main St^City^State^12345|(555)123-4567\rSPM|1|SPEC123^21321||BLD^Blood|||||||Routine|||Normal||||20241209113000|||||||\rOBR|1|ORD001||GLU^Glucose Test|||20241209110000||||||F|||||LAB123\rOBX|1|NM|GLU^Glucose Level||85|mg/dL|70-99|N|||F\rNTE|1|L|Specimens was hemolyzed; results may be affected.",
			},
			want: "OBX processed",
		},
		{
			name: "OUL_R22_With_ORC",
			fields: fields{
				AnalyzerUsecase: &analyzer.Usecase{
					ObservationRequestRepository: &observation_request.Repository{
						DB: db,
					},
					ObservationResultRepository: &observation_result.Repository{
						DB: db,
					},
				},
			},
			args: args{
				message: ex,
			},
			want: "OBX processed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HlSevenHandler{
				AnalyzerUsecase: tt.fields.AnalyzerUsecase,
			}
			got, err := h.HL7Handler(context.TODO(), tt.args.message)
			assert.Nil(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
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
