package subsystem

import (
	"github.com/Compogo/types/set"
)

const (
	Devices   Name = "devices"
	Hugetlb   Name = "hugetlb"
	Freezer   Name = "freezer"
	Pids      Name = "pids"
	NetCLS    Name = "net_cls"
	NetPrio   Name = "net_prio"
	PerfEvent Name = "perf_event"
	Cpuset    Name = "cpuset"
	Cpu       Name = "cpu"
	Cpuacct   Name = "cpuacct"
	Memory    Name = "memory"
	Blkio     Name = "blkio"
	Rdma      Name = "rdma"
	CGroupV2  Name = ""
)

type Name string

var AllNames = set.NewSet[Name](Devices, Hugetlb, Freezer, Pids, NetCLS, NetPrio, PerfEvent, Cpuset, Cpu, Cpuacct, Memory, Blkio, Rdma, CGroupV2)
