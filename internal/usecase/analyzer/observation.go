package analyzer

import (
	"context"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

func (u *Usecase) ProcessOULR22(ctx context.Context, data entity.OUL_R22) error {
	specimens := data.Specimens
	for i := range specimens {
		//if specimens[i].ObservationRequest != nil {
		//	err := u.ObservationRequestRepository.CreateMany(ctx, specimens[i].ObservationRequest)
		//	if err != nil {
		//		return err
		//	}
		//}
		//
		if specimens[i].ObservationResult != nil {
			err := u.ObservationResultRepository.CreateMany(ctx, specimens[i].ObservationResult)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
