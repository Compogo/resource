package domain

// Resource представляет информацию о ресурсе.
type Resource struct {
	Type  Type
	Limit uint64
	Usage uint64
}
