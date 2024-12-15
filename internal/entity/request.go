package entity

type GetManyRequest struct {
	Order  string   `query:"_order"`
	Sort   string   `query:"_sort"`
	Start  int      `query:"_start"`
	End    int      `query:"_end"`
	Search string   `query:"_search"`
	Query  string   `query:"q"`
	ID     []string `query:"id"`
}

func (g GetManyRequest) IsSortDesc() bool {
	if g.Order == "DESC" {
		return true
	}

	return false
}
