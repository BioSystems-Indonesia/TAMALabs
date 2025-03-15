package entity

type TestTypeSpecimen struct {
	TestTypeID   int64        `json:"test_type_id"`
	SpecimenType SpecimenType `json:"specimen_type"`
}
