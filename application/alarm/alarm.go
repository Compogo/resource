package alarm

import (
	"context"
	"sync/atomic"

	"github.com/Compogo/compogo/closer"
	"github.com/Compogo/resource/domain"
	"github.com/Compogo/resource/infrastructure/config/alarm"
	"github.com/Compogo/types/emitter"
)

type Alarm struct {
	OnAlarm emitter.Emitter[domain.State]

	config *alarm.Config
	state  atomic.Uint32

	closer closer.Closer
}

func NewAlarm(closer closer.Closer, config *alarm.Config) *Alarm {
	a := &Alarm{
		OnAlarm: emitter.NewEmitter[domain.State](),
		closer:  closer,
		config:  config,
	}

	a.state.Store(uint32(domain.Normal))

	return a
}

func (alarm *Alarm) State() domain.State {
	return domain.State(alarm.state.Load())
}

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
