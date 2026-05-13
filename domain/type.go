package domain

//go:generate stringer -type=Type

const (
	Memory Type = iota
	CPU
)

type Type uint8
