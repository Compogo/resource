package alarm

import (
	"github.com/Compogo/compogo"
	"github.com/Compogo/resource/application/manager"
	"github.com/Compogo/resource/infrastructure/config/alarm"
)

// Component — компонент алертов ресурсов.
// Подписывается на изменения ресурсов и генерирует события при достижении порогов.
var Component = &compogo.Component{
	Name: "resource.alarm",
	Dependencies: compogo.Components{
		&manager.Component,
		&alarm.Component,
	},
	Init: compogo.StepFunc(func(container compogo.Container) error {
		return container.Provide(NewAlarm)
	}),
	PreExecute: compogo.StepFunc(func(container compogo.Container) error {
		return container.Invoke(func(m *manager.Manager, a *Alarm) {
			m.OnChangeResource.Subscribe(a.OnChangeResource)
		})
	}),
}
