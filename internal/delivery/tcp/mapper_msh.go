package tcp

import (
	"time"

	"github.com/kardianos/hl7/h251"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

func getNamespaceID(field *h251.HD) string {
	if field == nil {
		return ""
	}
	return field.NamespaceID
}

func getMessageStructure(field h251.MSG) string {
	return field.MessageStructure
}

func getVersionID(field h251.VID) string {
	return field.VersionID
}

func getMessageDate(field h251.TS) string {
	return field.Format(time.RFC3339)
}

func mapMSHToEntity(msh *h251.MSH) entity.MSH {
	if msh == nil {
		return entity.MSH{}
	}

	sendingApplication := getNamespaceID(msh.SendingApplication)
	sendingFacility := getNamespaceID(msh.SendingFacility)
	receivingApplication := getNamespaceID(msh.ReceivingApplication)
	receivingFacility := getNamespaceID(msh.ReceivingFacility)
	messageType := getMessageStructure(msh.MessageType)
	version := getVersionID(msh.VersionID)
	messageDate := getMessageDate(msh.DateTimeOfMessage)

	return entity.MSH{
		MessageControlID:     msh.MessageControlID,
		SendingApplication:   sendingApplication,
		SendingFacility:      sendingFacility,
		ReceivingApplication: receivingApplication,
		ReceivingFacility:    receivingFacility,
		MessageType:          messageType,
		MessageVersion:       version,
		MessageDate:          messageDate,
	}
}
