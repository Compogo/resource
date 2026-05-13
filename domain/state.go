package domain

//go:generate stringer -type=State

const (
	Normal State = iota
	Warning
	Alarm
)

type State uint8
