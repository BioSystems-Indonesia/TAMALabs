package tcp

import (
	"fmt"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	"github.com/kardianos/hl7/h251"
)

func getObservationIdentifier(field h251.CE) (string, string) {
	return field.Identifier, field.Text
}

func getUnits(field *h251.CE) string {
	if field == nil {
		return ""
	}
	return field.Identifier
}

func mapObservationValueToValues(values []h251.VARIES) entity.JSONStringArray {
	if values == nil {
		return entity.JSONStringArray{}
	}

	var results entity.JSONStringArray
	for i := range values {
		results = append(results, fmt.Sprintf("%v", values[i]))
	}
	return results
}

func mapOBXToObservationResultEntity(obx *h251.OBX) entity.ObservationResult {
	if obx == nil {
		return entity.ObservationResult{}
	}

	testCode, description := getObservationIdentifier(obx.ObservationIdentifier)

	return entity.ObservationResult{
		TestCode:       testCode,
		Description:    description,
		Values:         mapObservationValueToValues(obx.ObservationValue),
		Type:           obx.ValueType,
		Unit:           getUnits(obx.Units),
		ReferenceRange: obx.ReferencesRange,
		Date:           obx.DateTimeOfTheObservation,
		AbnormalFlag:   obx.AbnormalFlags,
		Comments:       obx.ObservationResultStatus,
	}
}
