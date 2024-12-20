package rest

import (
	"fmt"
	"testing"

	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

func TestFeatureListHandler_filterTables(t *testing.T) {
	type args struct {
		feature entity.Tables
		req     entity.GetManyRequest
	}
	tests := []struct {
		name  string
		args  args
		want  entity.Tables
		want1 int
	}{
		{
			name: "success",
			args: args{
				feature: entity.TableObservationType,
				req: entity.GetManyRequest{
					Query: "Cal",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FeatureListHandler{}
			got, got1 := f.filterTables(tt.args.feature, tt.args.req)
			fmt.Println(got, got1)

		})
	}
}
