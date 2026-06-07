package pagination

type Params struct {
	Page   int
	Size   int
	Offset int
}

func NewParams(page, size int) Params {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 20
	}
	return Params{
		Page:   page,
		Size:   size,
		Offset: (page - 1) * size,
	}
}
