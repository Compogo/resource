package domain

type Resource struct {
	Type  Type
	Limit uint64
	Usage uint64
}
