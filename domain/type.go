package domain

//go:generate stringer -type=Type

const (
	// Memory оперативная память
	Memory Type = iota

	// CPU процессорное время
	CPU
)

// Type — тип ресурса.
type Type uint8
