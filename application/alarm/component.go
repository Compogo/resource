package alarm

import (
	"github.com/Compogo/compogo/component"
	"github.com/Compogo/compogo/container"
	"github.com/Compogo/resource/application/manager"
	"github.com/Compogo/resource/infrastructure/config/alarm"
)

var Component = &component.Component{
	Dependencies: component.Components{
		manager.Component,
		alarm.Component,
	},
	Init: component.StepFunc(func(container container.Container) error {
		return container.Provide(NewAlarm)
	}),
	PreExecute: component.StepFunc(func(container container.Container) error {
		return container.Invoke(func(m *manager.Manager, a *Alarm) {
			m.OnChangeResource.Subscribe(a.OnChangeResource)
		})
	}),
}
