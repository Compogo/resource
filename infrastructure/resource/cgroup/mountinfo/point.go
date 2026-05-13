package mountinfo

import (
	"github.com/Compogo/resource/infrastructure/resource/cgroup/path"
	"github.com/Compogo/types/set"
)

type Point struct {
	Id             uint32
	ParentID       uint32
	DeviceID       string
	Root           path.Path
	MountPoint     path.Path
	Options        set.Set[string]
	OptionalFields set.Set[string]
	FSType         string
	MountSource    string
	SuperOptions   set.Set[string]
}
