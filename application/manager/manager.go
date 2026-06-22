package manager

import (
	"context"
	"sync"

	"github.com/Compogo/compogo"
	"github.com/Compogo/resource/domain"
	"github.com/Compogo/resource/infrastructure/resource/cgroup"
	"github.com/Compogo/types/emitter"
)

// Manager — менеджер ресурсов.
// Собирает статистику использования CPU и Memory через cgroup.
// При изменении ресурсов генерирует событие OnChangeResource.
type Manager struct {
	OnChangeResource emitter.Emitter[*domain.Resource]

	cgroup *cgroup.CGroup

	rwm    sync.RWMutex
	cpu    domain.Resource
	memory domain.Resource

	closer compogo.Closer
}

// NewManager создаёт новый менеджер ресурсов.
func NewManager(closer compogo.Closer) *Manager {
	m := &Manager{
		OnChangeResource: emitter.NewEmitter[*domain.Resource](),
		cgroup:           cgroup.NewCGroup(),
		closer:           closer,
	}

	m.cpu.Type = domain.CPU
	m.memory.Type = domain.Memory

	return m
}

// Process собирает статистику ресурсов и генерирует события при изменениях.
// Вызывается периодически через Repeater.
func (manager *Manager) Process(_ context.Context) error {
	stat, err := manager.cgroup.Stat()
	if err != nil {
		return err
	}

	manager.rwm.Lock()
	defer manager.rwm.Unlock()

	emit := false

	if stat.CPU != nil && stat.CPU.Limit != manager.cpu.Limit {
		manager.cpu.Limit = stat.CPU.Limit
		emit = true
	}

	if stat.Cpuacct != nil && stat.Cpuacct.Usage != manager.cpu.Usage {
		manager.cpu.Usage = stat.Cpuacct.Usage
		emit = true
	}

	if emit {
		manager.OnChangeResource.Emit(manager.closer.GetContext(), &manager.cpu)
	}

	emit = false

	if stat.Memory != nil && stat.Memory.Limit != manager.memory.Limit {
		manager.memory.Limit = stat.Memory.Limit
		emit = true
	}

	if stat.Memory != nil && stat.Memory.Usage != manager.memory.Usage {
		manager.memory.Usage = stat.Memory.Usage
		emit = true
	}

	if emit {
		manager.OnChangeResource.Emit(manager.closer.GetContext(), &manager.memory)
	}

	return nil
}

// CPU возвращает текущий CPU ресурс.
func (manager *Manager) CPU() *domain.Resource {
	manager.rwm.RLock()
	defer manager.rwm.RUnlock()
	return &manager.cpu
}

// Memory возвращает текущий Memory ресурс.
func (manager *Manager) Memory() *domain.Resource {
	manager.rwm.RLock()
	defer manager.rwm.RUnlock()
	return &manager.memory
}
