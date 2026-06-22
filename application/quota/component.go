package quota

import (
	"github.com/Compogo/compogo"
	"github.com/Compogo/resource/application/manager"
)

// Component — компонент применения квот ресурсов.
// Подписывается на изменения ресурсов и применяет их к рантайму.
var Component = compogo.Component{
	Name: "resource.quota",
	Dependencies: compogo.Components{
		&manager.Component,
	},
	Init: compogo.StepFunc(func(container compogo.Container) error {
		return container.Provide(NewQuota)
	}),
	PreExecute: compogo.StepFunc(func(container compogo.Container) error {
		return container.Invoke(func(m *manager.Manager, q *Quota) {
			m.OnChangeResource.Subscribe(q.OnChangeResource)
		})
	}),
}
