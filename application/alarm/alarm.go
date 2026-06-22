package alarm

import (
	"context"
	"sync/atomic"

	"github.com/Compogo/compogo"
	"github.com/Compogo/resource/domain"
	"github.com/Compogo/resource/infrastructure/config/alarm"
	"github.com/Compogo/types/emitter"
)

// Alarm — система алертов для ресурсов.
// Отслеживает использование ресурсов и генерирует события при достижении порогов.
type Alarm struct {
	OnAlarm emitter.Emitter[domain.State]

	config *alarm.Config
	state  atomic.Uint32

	closer compogo.Closer
}

// NewAlarm создаёт новый Alarm.
func NewAlarm(closer compogo.Closer, config *alarm.Config) *Alarm {
	a := &Alarm{
		OnAlarm: emitter.NewEmitter[domain.State](),
		closer:  closer,
		config:  config,
	}

	a.state.Store(uint32(domain.Normal))

	return a
}

// State возвращает текущее состояние.
func (alarm *Alarm) State() domain.State {
	return domain.State(alarm.state.Load())
}

// OnChangeResource вызывается при изменении ресурсов.
// Вычисляет процент использования и обновляет состояние.
func (alarm *Alarm) OnChangeResource(_ context.Context, resource *domain.Resource) {
	resourcePercentUsage := float32(float64(resource.Usage) / float64(resource.Limit))
	newState := domain.Normal

	if resourcePercentUsage >= alarm.config.WarnStatePercent {
		newState = domain.Warning
	}

	if resourcePercentUsage >= alarm.config.AlarmStatePercent {
		newState = domain.Alarm
	}

	oldState := domain.State(alarm.state.Swap(uint32(newState)))

	if oldState != newState {
		alarm.OnAlarm.Emit(alarm.closer.GetContext(), newState)
	}
}
