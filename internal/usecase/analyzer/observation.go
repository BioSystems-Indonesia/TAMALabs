package analyzer

import (
	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

func (u *Usecase) ProcessOULR22(data entity.OUL_R22) error {
	err := u.ObservationRequestRepository.Create(&data.ObservationRequest)
	if err != nil {
		return err
	}

	err = u.ObservationRepository.CreateMany(&data.Observations)
	if err != nil {
		return err
	}

	err = u.SpecimenRepository.Create(&data.Specimen)
	if err != nil {
		return err
	}
	return nil
}
