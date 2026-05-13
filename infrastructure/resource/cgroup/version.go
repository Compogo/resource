package cgroup

//go:generate stringer --type=Version

const (
	CGroupV1 Version = iota
	CGroupV2
)

type Version uint8
