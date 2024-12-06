package hl_seven

import "github.com/oibacidem/lims-hl-seven/internal/entity"

// SendORM returns an ORM struct
func (u *Usecase) SendORM(message entity.SendORMRequest) (*entity.SendORMResponse, error) {
	resp, err := u.ORMRepository.SendORM(message)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
