package pg

type Pagination struct {
	Cursor   int
	PageSize int
}

func PaginationSettings(cursor, pageSize *int) Pagination {
	p := Pagination{}

	if cursor != nil {
		p.Cursor = *cursor
	} else {
		p.Cursor = 0
	}

	if pageSize != nil {
		p.PageSize = *pageSize
	} else {
		p.PageSize = 100
	}

	return p
}
