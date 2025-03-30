package tcp

import (
	"fmt"

	"github.com/kardianos/hl7/h251"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

func MapQBPQ11ToEntity(msg *h251.QBP_Q11) (entity.QBP_Q11, error) {
	msh := mapMSHToEntity(msg.MSH)

	return entity.QBP_Q11{
		Msh:     msh,
		Barcode: fmt.Sprintf("%v", msg.QPD.UserParametersInSuccessiveFields),
	}, nil
}
