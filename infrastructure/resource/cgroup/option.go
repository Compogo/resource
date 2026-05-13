package cgroup

import (
	"github.com/Compogo/resource/infrastructure/resource/cgroup/mountinfo"
	"github.com/Compogo/resource/infrastructure/resource/cgroup/subsystem"
)

type Option func(cgroup *CGroup)

func WithSubsystemParser(subsystem *subsystem.Parser) Option {
	return func(cgroup *CGroup) {
		cgroup.subsystemParser = subsystem
	}
}

func WithMountinfoParser(mountinfo *mountinfo.Parser) Option {
	return func(cgroup *CGroup) {
		cgroup.mountinfoParser = mountinfo
	}
}

func WithPID(pid uint64) Option {
	return func(cgroup *CGroup) {
		cgroup.subsystemParser = subsystem.NewParser(subsystem.WithPath(subsystem.PidPath(pid)))
		cgroup.mountinfoParser = mountinfo.NewParser(mountinfo.WithPath(mountinfo.PidPath(pid)))
	}
}

func WithVersionOrder(versions ...Version) Option {
	return func(cgroup *CGroup) {
		cgroup.order = versions
	}
}

func WithOnlyVersion(versions ...Version) Option {
	return func(cgroup *CGroup) {
		cgroup.onlyVersion.Reset()
		cgroup.onlyVersion.Add(versions...)
	}
}
