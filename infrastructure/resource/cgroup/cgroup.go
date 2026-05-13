package cgroup

import (
	"github.com/Compogo/resource/infrastructure/resource/cgroup/mountinfo"
	"github.com/Compogo/resource/infrastructure/resource/cgroup/subsystem"
	"github.com/Compogo/types/linker"
	"github.com/Compogo/types/set"
)

type Loader func(*linker.Linker[subsystem.Name, *subsystem.Subsystem], set.Set[*mountinfo.Point], *Stat) (*Stat, error)

var loaders = linker.NewLinker[Version, Loader](
	linker.Link(CGroupV1, Loader(LoaderV1)),
	linker.Link(CGroupV2, Loader(LoaderV2)),
)

type CGroup struct {
	subsystemParser *subsystem.Parser
	mountinfoParser *mountinfo.Parser

	onlyVersion set.Set[Version]
	order       []Version
}

func NewCGroup(options ...Option) *CGroup {
	cgroup := &CGroup{}

	options = append([]Option{
		WithSubsystemParser(subsystem.NewParser()),
		WithMountinfoParser(mountinfo.NewParser()),
		WithOnlyVersion(CGroupV1, CGroupV2),
	}, options...)

	return cgroup
}

func (cgroup *CGroup) Stat() (*Stat, error) {
	subsystems, err := cgroup.subsystemParser.Subsystems()
	if err != nil {
		return nil, err
	}

	points, err := cgroup.mountinfoParser.Points()
	if err != nil {
		return nil, err
	}

	stat := &Stat{}
	var loader Loader
	for _, version := range cgroup.order {
		if cgroup.onlyVersion.Len() > 0 && !cgroup.onlyVersion.Contains(version) {
			continue
		}

		loader, err = loaders.Get(version)
		if err != nil {
			return nil, err
		}

		stat, err = loader(subsystems, points, stat)
		if err != nil {
			return nil, err
		}
	}

	return stat, nil
}
