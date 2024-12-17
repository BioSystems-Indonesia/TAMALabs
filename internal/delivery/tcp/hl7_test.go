package tcp

import (
	"github.com/oibacidem/lims-hl-seven/internal/usecase/analyzer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHlSevenHandler(t *testing.T) {
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
			name: "OUL_R22_ELGATAMA",
			fields: fields{
				AnalyzerUsecase: &analyzer.Usecase{},
			},
			args: args{
				message: "MSH|^~\\&|BA200|Biosystems|Host|Host provider|20241209214354||OUL^R22^OUL_R22|fac365d1-a158-4ce0-97f9-95afa1359571|P|2.5.1|||ER|AL||UNICODE UTF-8|||LAB-29^IH\rOBX|1|NM|UREA-BUN-UV^UREA-BUN-UV^A400||60.7323189|mg/dL^mg/dL^A400||NONE|||F|||||ADMIN||A400^Biosystems~834000237^Biosystems|20241210044024",
			},
			want: "OBX processed",
		},
		{
			name: "OUL_R22",
			fields: fields{
				AnalyzerUsecase: &analyzer.Usecase{},
			},
			args: args{
				message: "MSH|^~\\&|LabSystem|LabFacility|EHR|Hospital|20241209120000||OUL^R22|MSG001|P|2.5.1\rPID|12|123|123456||John^Doe||19800101|M|||123 Main St^City^State^12345|(555)123-4567\rSPM|1|SPEC123||BLD^Blood|||||||Routine|||Normal||||20241209113000|||||||\rOBR|1|ORD001||GLU^Glucose Test|||20241209110000||||||F|||||LAB123\rOBX|1|NM|GLU^Glucose Level||85|mg/dL|70-99|N|||F\rNTE|1|L|Specimen was hemolyzed; results may be affected.",
			},
			want: "OBX processed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HlSevenHandler{
				AnalyzerUsecase: tt.fields.AnalyzerUsecase,
			}
			got, err := h.HL7Handler(tt.args.message)
			assert.Nil(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
