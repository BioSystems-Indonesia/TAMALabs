package tcp

import (
	"fmt"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	"github.com/kardianos/hl7/h251"
)

func MapQBPQ11ToEntity(msg *h251.QBP_Q11) (entity.QBP_Q11, error) {
	msh := mapMSHToEntity(msg.MSH)

	return entity.QBP_Q11{
		Msh: msh,
		QPD: entity.QPD{
			QueryTag: msg.QPD.QueryTag,
			Barcode: func() string {
				if msg.QPD.UserParametersInSuccessiveFields == nil {
					return ""
				}
				return fmt.Sprintf("%v", *msg.QPD.UserParametersInSuccessiveFields)
			}(),
		},
	}, nil
}
