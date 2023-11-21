package pagination

type Pagination interface {
	GetPageNumber() int32
	GetShowNumber() int32
}
