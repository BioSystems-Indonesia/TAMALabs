package analyzer

import (
	"context"
	"log"

	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

func (u *Usecase) ProcessOULR22(ctx context.Context, data entity.OUL_R22) error {
	specimens := data.Specimens
	for i := range specimens {
		spEntities, err := u.SpecimenRepository.FindByBarcode(ctx, specimens[i].HL7ID)
		if err != nil {
			log.Println(err)
			continue
		}
		var spEntity entity.Specimen

		if len(spEntities) > 0 {
			spEntity = spEntities[0]
		}

		//if specimens[i].ObservationRequest != nil {
		//	err := u.ObservationRequestRepository.CreateMany(ctx, specimens[i].ObservationRequest)
		//	if err != nil {
		//		return err
		//	}
		//}
		//
		if specimens[i].ObservationResult != nil {
			for j := range specimens[i].ObservationResult {
				specimens[i].ObservationResult[j].SpecimenID = int64(spEntity.ID)
			}
			err := u.ObservationResultRepository.CreateMany(ctx, specimens[i].ObservationResult)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
