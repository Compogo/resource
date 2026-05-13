package cgroup

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Compogo/resource/infrastructure/resource/cgroup/mountinfo"
	"github.com/Compogo/resource/infrastructure/resource/cgroup/path"
	"github.com/Compogo/resource/infrastructure/resource/cgroup/subsystem"
	"github.com/Compogo/types/linker"
	"github.com/Compogo/types/mapper"
	"github.com/Compogo/types/set"
	"github.com/spf13/cast"
)

const (
	v2MountPoints = "/sys/fs/cgroup"

	v2CpuMaxFileName  = string(subsystem.Cpu) + ".max"
	v2CpuStatFileName = string(subsystem.Cpu) + ".stat"

	v2MemoryMaxFileName     = string(subsystem.Memory) + ".max"
	v2MemoryCurrentFileName = string(subsystem.Memory) + ".current"

	v2CpuMinFields        = 2
	v2CpuLimitFieldIndex  = 0
	v2CpuPeriodFieldIndex = 1

	v2CpuStatMinFields       = 2
	v2CpuStatUsageFieldName  = "usage_usec"
	v2CpuStatFieldNameIndex  = 0
	v2CpuStatFieldValueIndex = 1
)

func LoaderV2(subsystems *linker.Linker[subsystem.Name, *subsystem.Subsystem], _ set.Set[*mountinfo.Point], stat *Stat) (*Stat, error) {
	v2Subsystem, err := subsystems.Get(subsystem.CGroupV2)
	if errors.Is(err, mapper.DoesNotExistError) {
		return nil, nil
	}

	stat, err = loaderV2(v2Subsystem.MountPointRootPath, stat)
	if err != nil {
		return nil, err
	}

	return stat, nil
}

func loaderV2(cgroupPath path.Path, stat *Stat) (*Stat, error) {
	var err error

	basePath := filepath.Join(v2MountPoints, string(cgroupPath))

	if stat.Memory == nil {
		stat, err = loaderV2Memory(basePath, stat)
		if err != nil {
			return nil, err
		}
	}

	if stat.CPU == nil {
		stat, err = loaderV2Cpu(basePath, stat)
		if err != nil {
			return nil, err
		}
	}

	if stat.Cpuacct == nil {
		stat, err = loaderV2Cpuacct(basePath, stat)
		if err != nil {
			return nil, err
		}
	}

	return stat, nil
}

func loaderV2Memory(cgroupPath string, stat *Stat) (*Stat, error) {
	limit, err := subsystem.ReadParamUint64(filepath.Join(cgroupPath, v2MemoryMaxFileName))
	if err != nil {
		return nil, err
	}

	if limit == 0 {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		limit = m.Sys
	}

	usage, err := subsystem.ReadParamUint64(filepath.Join(cgroupPath, v2MemoryCurrentFileName))
	if err != nil {
		return nil, err
	}

	stat.Memory = &Memory{
		Version: CGroupV2,
		Limit:   limit,
		Usage:   usage,
	}

	return stat, nil
}

func loaderV2Cpu(cgroupPath string, stat *Stat) (*Stat, error) {
	content, err := subsystem.ReadParam(filepath.Join(cgroupPath, v2CpuMaxFileName))
	if err != nil {
		return nil, err
	}

	if content == "" {
		return stat, nil
	}

	fields := strings.Fields(content)
	if len(fields) < v2CpuMinFields {
		return nil, fmt.Errorf("invalid cpu.max format: %s", content)
	}

	var limit uint64
	if fields[v2CpuLimitFieldIndex] == subsystem.MaxTextValue {
		limit = 0
	} else {
		limit, err = cast.ToUint64E(fields[v2CpuLimitFieldIndex])
		if err != nil {
			return nil, err
		}
	}

	period := DefaultPeriodCPU
	if len(fields) >= v2CpuMinFields && fields[v2CpuPeriodFieldIndex] != "" {
		period, err = cast.ToUint64E(fields[v2CpuPeriodFieldIndex])
		if err != nil {
			return nil, err
		}
	}

	if limit == 0 {
		limit = uint64(runtime.NumCPU()) * period
	}

	stat.CPU = &CPU{
		Version: CGroupV2,
		Limit:   limit,
		Period:  period,
	}

	return stat, nil
}

func loaderV2Cpuacct(cgroupPath string, stat *Stat) (*Stat, error) {
	cpuStatFile, err := os.Open(filepath.Join(cgroupPath, v2CpuStatFileName))
	if err != nil && os.IsNotExist(err) {
		return stat, nil
	}

	if err != nil {
		return nil, err
	}

	defer cpuStatFile.Close()

	lineScanner := bufio.NewScanner(cpuStatFile)
	lineScanner.Split(bufio.ScanLines)

	for lineScanner.Scan() {
		fields := strings.Fields(lineScanner.Text())

		if len(fields) != v2CpuStatMinFields {
			continue
		}

		if fields[v2CpuStatFieldNameIndex] == v2CpuStatUsageFieldName {
			usage, err := cast.ToUint64E(fields[v2CpuStatFieldValueIndex])
			if err != nil {
				return nil, err
			}

			stat.Cpuacct = &Cpuacct{
				Version: CGroupV2,
				Usage:   usage,
			}

			break
		}
	}

	if err := lineScanner.Err(); err != nil {
		return nil, err
	}

	return stat, nil
}
