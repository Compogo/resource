package subsystem

import "github.com/Compogo/resource/infrastructure/resource/cgroup/path"

type Subsystem struct {
	Id                 uint8
	Name               Name
	MountPointRootPath path.Path
}
