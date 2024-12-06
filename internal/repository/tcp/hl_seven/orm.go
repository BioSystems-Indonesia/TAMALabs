package hl_seven

import "github.com/oibacidem/lims-hl-seven/internal/entity"

// SendORM sends an ORM request to the HL7 server
func (r *Repository) SendORM(request entity.SendORMRequest) (*entity.SendORMResponse, error) {
	message := request.ORM.Serialize()

	// send message to tcp
	resp, err := r.tcp.Send(message)
	if err != nil {
		return nil, err
	}

	respParsed, err := entity.ACKMessage(resp).Parse()
	if err != nil {
		return nil, err
	}
	return &entity.SendORMResponse{
		ACK: respParsed,
	}, nil
}
