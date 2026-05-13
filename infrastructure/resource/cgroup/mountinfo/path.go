package mountinfo

import (
	"fmt"

	"github.com/Compogo/resource/infrastructure/resource/cgroup/path"
)

func PidPath(pid uint64) path.Path {
	return path.Path(fmt.Sprintf(path.CGroupMountInfoPidPathFormat, pid))
}
