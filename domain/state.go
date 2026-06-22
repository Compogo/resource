package domain

//go:generate stringer -type=State

const (
	// Normal нормальное состояние
	Normal State = iota

	// Warning предупреждение (используется >70%)
	Warning

	// Alarm тревога (используется >85%)
	Alarm
)

// State — состояние ресурса.
type State uint8
