package quota

import (
	"context"
	"runtime"
	"runtime/debug"

	"github.com/Compogo/resource/domain"
	"github.com/Compogo/resource/infrastructure/resource/cgroup"
)

// Quota — применяет ограничения ресурсов к рантайму Go.
// При изменении лимитов CPU/Memory через cgroup, автоматически настраивает:
//   - GOMAXPROCS — количество процессоров
//   - SetMemoryLimit — лимит памяти для GC
type Quota struct{}

// NewQuota создаёт новый Quota.
func NewQuota() *Quota {
	return &Quota{}
}

// OnChangeResource вызывается при изменении ресурсов.
// Применяет новые лимиты к рантайму.
func (quota *Quota) OnChangeResource(_ context.Context, resource *domain.Resource) {
	switch resource.Type {
	case domain.CPU:
		runtime.GOMAXPROCS(int(resource.Limit / cgroup.DefaultPeriodCPU))
	case domain.Memory:
		debug.SetMemoryLimit(int64(resource.Limit))
	}
}
