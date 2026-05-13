package manager

import (
	"github.com/Compogo/compogo/component"
	"github.com/Compogo/compogo/container"
	"github.com/Compogo/repeater"
	"github.com/Compogo/resource/infrastructure/config/manager"
)

var Component = &component.Component{
	Name: "resource.manager",
	Dependencies: component.Components{
		manager.Component,
		repeater.Component,
	},
	Init: component.StepFunc(func(container container.Container) error {
		return container.Provide(NewManager)
	}),
	PreExecute: component.StepFunc(func(container container.Container) error {
		return container.Invoke(func(r repeater.Repeater, config *manager.Config, m *Manager) error {
			return r.AddProcess(repeater.NewTask("resource.manager", config.Delay, m.Process))
		})
	}),
}
