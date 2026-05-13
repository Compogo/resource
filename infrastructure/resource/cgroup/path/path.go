package path

const (
	CGroupSubsystemSelfPath      Path = "/proc/self/cgroup"
	CGroupSubsystemPidPathFormat      = "/proc/%d/cgroup"

	CGroupMountInfoSelfPath      Path = "/proc/self/mountinfo"
	CGroupMountInfoPidPathFormat      = "/proc/%d/mountinfo"
)

type Path string

func MountPath(p string) Path {
	return Path(p)
}
