package cgroup

import (
	"errors"
	"path/filepath"
	"runtime"

	"github.com/Compogo/resource/infrastructure/resource/cgroup/mountinfo"
	"github.com/Compogo/resource/infrastructure/resource/cgroup/subsystem"
	"github.com/Compogo/types/linker"
	"github.com/Compogo/types/mapper"
	"github.com/Compogo/types/set"
)

const (
	v1CpuQuotaFileName  = string(subsystem.Cpu) + ".cfs_quota_us"
	v1CpuPeriodFileName = string(subsystem.Cpu) + ".cfs_period_us"
	v1CpuUsageFileName  = string(subsystem.Cpuacct) + ".usage"

	v1MemoryQuotaFileName = string(subsystem.Memory) + ".limit_in_bytes"
	v1MemoryUsageFileName = string(subsystem.Memory) + ".usage_in_bytes"
)

func LoaderV1(subsystems *linker.Linker[subsystem.Name, *subsystem.Subsystem], mountPoints set.Set[*mountinfo.Point], stat *Stat) (*Stat, error) {
	var err error

	if stat.Memory == nil {
		stat, err = loaderV1Memory(subsystems, mountPoints, stat)
		if err != nil {
			return nil, err
		}
	}

	if stat.CPU == nil {
		stat, err = loaderV1Cpu(subsystems, mountPoints, stat)
		if err != nil {
			return nil, err
		}
	}

	if stat.Cpuacct == nil {
		stat, err = loaderV1Cpuacct(subsystems, mountPoints, stat)
		if err != nil {
			return nil, err
		}
	}

	return stat, nil
}

func loaderV1Memory(subsystems *linker.Linker[subsystem.Name, *subsystem.Subsystem], mountPoints set.Set[*mountinfo.Point], stat *Stat) (*Stat, error) {
	memSubsystem, err := subsystems.Get(subsystem.Memory)
	if errors.Is(err, mapper.DoesNotExistError) {
		return stat, nil
	}

	for point := range mountPoints {
		if point.Root != memSubsystem.MountPointRootPath {
			continue
		}

		if !point.SuperOptions.Contains(string(subsystem.Memory)) {
			continue
		}

		basePath := filepath.Join(string(point.MountPoint), string(memSubsystem.MountPointRootPath))

		limit, err := subsystem.ReadParamUint64(filepath.Join(basePath, v1MemoryQuotaFileName))
		if err != nil {
			return nil, err
		}

		if limit == 0 {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)

			limit = m.Sys
		}

		usage, err := subsystem.ReadParamUint64(filepath.Join(basePath, v1MemoryUsageFileName))
		if err != nil {
			return nil, err
		}

		stat.Memory = &Memory{
			Version: CGroupV1,
			Limit:   limit,
			Usage:   usage,
		}

		break
	}

	return stat, nil
}

func loaderV1Cpu(subsystems *linker.Linker[subsystem.Name, *subsystem.Subsystem], mountPoints set.Set[*mountinfo.Point], stat *Stat) (*Stat, error) {
	cpuSubsystem, err := subsystems.Get(subsystem.Cpu)
	if errors.Is(err, mapper.DoesNotExistError) {
		return stat, nil
	}

	for point := range mountPoints {
		if point.Root != cpuSubsystem.MountPointRootPath {
			continue
		}

		if !point.SuperOptions.Contains(string(subsystem.Cpu)) {
			continue
		}

		basePath := filepath.Join(string(point.MountPoint), string(cpuSubsystem.MountPointRootPath))

		period, err := subsystem.ReadParamUint64(filepath.Join(basePath, v1CpuPeriodFileName))
		if err != nil {
			return nil, err
		}
		if period == 0 {
			period = DefaultPeriodCPU
		}

		limit, err := subsystem.ReadParamInt64(filepath.Join(basePath, v1CpuQuotaFileName))
		if err != nil {
			return nil, err
		}

		if limit <= 0 {
			limit = int64(uint64(runtime.NumCPU()) * period)
		}

		stat.CPU = &CPU{
			Version: CGroupV1,
			Limit:   uint64(limit),
			Period:  period,
		}

		break
	}

	return stat, nil
}

func loaderV1Cpuacct(subsystems *linker.Linker[subsystem.Name, *subsystem.Subsystem], mountPoints set.Set[*mountinfo.Point], stat *Stat) (*Stat, error) {
	var point *mountinfo.Point

	cpuSubsystem, err := subsystems.Get(subsystem.Cpuacct)
	if errors.Is(err, mapper.DoesNotExistError) {
		cpuSubsystem, err = subsystems.Get(subsystem.Cpu)
		if errors.Is(err, mapper.DoesNotExistError) {
			return stat, nil
		}

		for findPoint := range mountPoints {
			if findPoint.Root != cpuSubsystem.MountPointRootPath {
				continue
			}

			if !findPoint.SuperOptions.Contains(string(subsystem.Cpu)) {
				continue
			}

			point = findPoint
			break
		}
	} else {
		for findPoint := range mountPoints {
			if findPoint.Root != cpuSubsystem.MountPointRootPath {
				continue
			}

			if !findPoint.SuperOptions.Contains(string(subsystem.Cpuacct)) {
				continue
			}

			point = findPoint
			break
		}
	}

	if point != nil {
		basePath := filepath.Join(string(point.MountPoint), string(cpuSubsystem.MountPointRootPath))

		usage, err := subsystem.ReadParamUint64(filepath.Join(basePath, v1CpuUsageFileName))
		if err != nil {
			return nil, err
		}

		stat.Cpuacct = &Cpuacct{
			Version: CGroupV1,
			Usage:   usage,
		}
	}

	return stat, nil
}
