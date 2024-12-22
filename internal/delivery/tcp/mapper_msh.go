package tcp

import (
	"github.com/kardianos/hl7/h251"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"time"
)

func mapMSHToEntity(msh *h251.MSH) entity.MSH {
	return entity.MSH{
		MessageControlID:     msh.MessageControlID,
		SendingApplication:   msh.SendingApplication.NamespaceID,
		SendingFacility:      msh.SendingFacility.NamespaceID,
		ReceivingApplication: msh.ReceivingApplication.NamespaceID,
		ReceivingFacility:    msh.ReceivingFacility.NamespaceID,
		MessageType:          msh.MessageType.MessageStructure,
		MessageVersion:       msh.VersionID.VersionID,
		MessageDate:          msh.DateTimeOfMessage.Format(time.RFC3339),
	}
}
