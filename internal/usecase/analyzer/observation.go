package analyzer

import (
	"fmt"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

func (u *Usecase) ProcessOULR22(data entity.OUL_R22) error {
	fmt.Println(data)
	return nil
}
