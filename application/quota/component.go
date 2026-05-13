package quota

import (
	"github.com/Compogo/compogo/component"
	"github.com/Compogo/compogo/container"
	"github.com/Compogo/resource/application/manager"
)

var Component = &component.Component{
	Dependencies: component.Components{
		manager.Component,
	},
	Init: component.StepFunc(func(container container.Container) error {
		return container.Provide(NewQuota)
	}),
	PreExecute: component.StepFunc(func(container container.Container) error {
		return container.Invoke(func(m *manager.Manager, q *Quota) {
			m.OnChangeResource.Subscribe(q.OnChangeResource)
		})
	}),
}
