package cgroup

const (
	DefaultPeriodCPU = uint64(100000)
)

type Stat struct {
	CPU     *CPU
	Cpuacct *Cpuacct
	Memory  *Memory
}

type CPU struct {
	Version Version
	Limit   uint64
	Period  uint64
}

type Cpuacct struct {
	Version Version
	Usage   uint64
}

type Memory struct {
	Version Version
	Limit   uint64
	Usage   uint64
}
