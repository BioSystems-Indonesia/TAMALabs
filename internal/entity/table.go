package entity

import (
	"slices"
	"strings"
)

type Table struct {
	ID             string      `json:"id"`
	Name           string      `json:"name"`
	AdditionalInfo interface{} `json:"additional_info"`
}

type Tables []Table

func (t Tables) Find(id string) (Table, bool) {
	for _, table := range t {
		if table.ID == id {
			return table, true
		}
	}
	return Table{}, false
}

func (t Tables) FilterID(id []string) []Table {
	var tables []Table
	for _, table := range t {
		if slices.Contains(id, table.ID) {
			tables = append(tables, table)
		}
	}
	return tables
}

func (t Tables) FilterName(search string) []Table {
	var tables []Table
	for _, table := range t {
		if strings.Contains(strings.ToLower(table.Name), strings.ToLower(search)) {
			tables = append(tables, table)
		}
	}
	return tables
}

var TableSpecimenType = Tables{
	{
		ID: string(SpecimenTypeSER), Name: `SER`,
	},
	{
		ID: string(SpecimenTypeURI), Name: `URI`,
	},
	{
		ID: string(SpecimenTypeCSF), Name: `CSF`,
	},
	{
		ID: string(SpecimenTypeLIQ), Name: `LIQ`,
	},
	{
		ID: string(SpecimenTypePLM), Name: `PLM`,
	},
	{
		ID: string(SpecimenTypeSEM), Name: `SEM`,
	},
	{
		ID: string(SpecimenTypeWBL), Name: `WBL`,
	},
}

var TableSex = Tables{
	{
		ID: `M`, Name: `Male`,
	},
	{
		ID: `F`, Name: `Female`,
	},
}

var TableSpecimenTest = Tables{
	{
		ID: `CHOLESTEROL`, Name: `CHOLESTEROL`,
	},
	{
		ID: `BUN`, Name: `BUN`,
	},
}

var TableWorkOrderStatus = Tables{
	{
		ID: string(WorkOrderStatusNew), Name: `NEW`,
	},
	{
		ID: string(WorkOrderCancelled), Name: `CANCELLED`,
	},
	{
		ID: string(WorkOrderStatusPending), Name: `PENDING`,
	},
	{
		ID: string(WorkOrderStatusCompleted), Name: `COMPLETED`,
	},
}

var TableBaudRate = Tables{
	{ID: "9600", Name: "9600"},
	{ID: "19200", Name: "19200"},
	{ID: "38400", Name: "38400"},
	{ID: "57600", Name: "57600"},
	{ID: "115200", Name: "115200"},
	{ID: "230400", Name: "230400"},
	{ID: "460800", Name: "460800"},
	{ID: "921600", Name: "921600"},
}

type ObservationInfo struct {
	Type ObservationTestType `json:"type"`
}

type ObservationTestType string

const (
	ObservationTestTypeStandard ObservationTestType = "STANDARD_TEST"
)

var TableList = map[string]Tables{
	"sex":               TableSex,
	"specimen-type":     TableSpecimenType,
	"specimen-test":     TableSpecimenTest,
	"work-order-status": TableWorkOrderStatus,
	"device-type":       TableDeviceType,
	"baud-rate":         TableBaudRate,
}
