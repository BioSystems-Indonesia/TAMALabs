package khanzauc

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestUseCase_ProcessRequest(t *testing.T) {
	t.Skip()

	u := &Usecase{}

	err := u.ProcessRequest(context.Background(), []byte(rawRequest))

	assert.NoError(t, err)

}

func TestUseCase_UnmarshalRequest(t *testing.T) {
	var r Request

	err := json.Unmarshal([]byte(rawRequest), &r)

	assert.NoError(t, err)

	u := &Usecase{}

	p, err := u.convertIntoPatient(r)
	assert.NoError(t, err)

	assert.Equal(t, "Pasien", p.FirstName)
	assert.Equal(t, "Pertama", p.LastName)
	assert.Equal(t, entity.PatientSexMale, p.Sex)
	// TODO should we use UTC or Local for Birthdate
	//assert.Equal(t, time.Date(2020, 1, 2, 0, 0, 0, 0, time.Local), p.Birthdate)
	assert.Equal(t, time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC), p.Birthdate)
}

const rawRequest = `
{ 
	"order": {
		"msh": {
			"product": "product_soft_medix",
			"version": "version_of_softmedix",
			"user_id": "user_of_soft_medix",
			"key": "somekey" 
		}, 
		"pid": {
			"pmrn": "000051",
			"pname": "Pasien Pertama",
			"sex": "L",
			"birth_dt": "2020-01-02",
			"address": "",
			"no_tlp": "",
			"no_hp": "",
			"email": "",
			"nik": "nik"
		},
		"obr": {
			"order_control": "N",
			"ptype": "OP",
			"reg_no": "",
			"order_lab": "",
			"provider_id": "",
			"provider_name": "",
			"order_date": "",
			"clinician_id": "",
			"clinician_name": "",
			"bangsal_id": "",
			"bangsal_name": "",
			"bed_id": "",
			"bed_name": "",
			"class_id": "",
			"class_name": "",
			"cito": "N",
			"med_legal": "N",
			"user_id": "",
			"reserve1": "",
			"reserve2": "",
			"reserve3": "",
			"reserve4": "",
			"order_test": ["satu", "dua"]
		}
	}
}
`
