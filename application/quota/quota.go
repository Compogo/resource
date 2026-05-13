package quota

import (
	"context"
	"runtime"
	"runtime/debug"

	"github.com/Compogo/resource/domain"
	"github.com/Compogo/resource/infrastructure/resource/cgroup"
)

type Quota struct{}

func NewQuota() *Quota {
	return &Quota{}
}

func (quota *Quota) OnChangeResource(_ context.Context, resource *domain.Resource) {
	switch resource.Type {
	case domain.CPU:
		runtime.GOMAXPROCS(int(resource.Limit / cgroup.DefaultPeriodCPU))
	case domain.Memory:
		debug.SetMemoryLimit(int64(resource.Limit))
	}
}
