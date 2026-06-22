package manager

import (
	"github.com/Compogo/compogo"
	"github.com/Compogo/repeater"
	"github.com/Compogo/resource/infrastructure/config/manager"
)

// Component — компонент менеджера ресурсов.
// Запускает периодический сбор статистики через Repeater.
var Component = compogo.Component{
	Name: "resource.manager",
	Dependencies: compogo.Components{
		&manager.Component,
		&repeater.Component,
	},
	Init: compogo.StepFunc(func(container compogo.Container) error {
		return container.Provide(NewManager)
	}),
	Execute: compogo.StepFunc(func(container compogo.Container) error {
		return container.Invoke(func(r repeater.Repeater, config *manager.Config, m *Manager) error {
			return r.AddProcess(repeater.NewTask("resource.manager", config.Delay, m.Process))
		})
	}),
}
