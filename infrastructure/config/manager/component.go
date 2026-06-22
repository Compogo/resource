package manager

import (
	"github.com/Compogo/compogo"
	"github.com/Compogo/compogo/flag"
)

// Component — компонент конфигурации менеджера ресурсов.
var Component = compogo.Component{
	Name: "resource.manager.config",
	Init: compogo.StepFunc(func(container compogo.Container) error {
		return container.Provide(NewConfig)
	}),
	BindFlags: compogo.BindFlags(func(flagSet flag.FlagSet, container compogo.Container) error {
		return container.Invoke(func(config *Config) {
			flagSet.DurationVar(&config.Delay, DelayFieldName, DelayDefault, "frequency of statistics collection")
		})
	}),
	Configuration: compogo.StepFunc(func(container compogo.Container) error {
		return container.Invoke(Configuration)
	}),
}
