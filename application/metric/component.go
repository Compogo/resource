package metric

import (
	"github.com/Compogo/compogo"
	"github.com/Compogo/resource/application/manager"
)

// Component — компонент метрик ресурсов.
// Подписывается на изменения ресурсов и обновляет метрики.
var Component = compogo.Component{
	Name: "resource.metrics",
	Dependencies: compogo.Components{
		&manager.Component,
	},
	Init: compogo.StepFunc(func(container compogo.Container) error {
		return container.Provide(NewMetric)
	}),
	PreExecute: compogo.StepFunc(func(container compogo.Container) error {
		return container.Invoke(func(m *manager.Manager, metric *Metric) {
			m.OnChangeResource.Subscribe(metric.OnChangeResource)
		})
	}),
}
