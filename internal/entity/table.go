package entity

import (
	"slices"
	"strings"
)

type Table struct {
	ID   string `json:"id"`
	Name string `json:"name"`
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
		if strings.Contains(table.Name, search) {
			tables = append(tables, table)
		}
	}
	return tables
}

var TableSpecimentType = Tables{
	{
		ID: `SER`, Name: `Serum`,
	},
}

var TableSex = Tables{
	{
		ID: `M`, Name: `Male`,
	},
	{
		ID: `F`, Name: `Female`,
	},
	{
		ID: `U`, Name: `Unknown`,
	},
}

var TableSpecimentTest = Tables{
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

var TableList = map[string]Tables{
	"sex":               TableSex,
	"speciment-type":    TableSpecimentType,
	"speciment-test":    TableSpecimentTest,
	"work-order-status": TableWorkOrderStatus,
}
