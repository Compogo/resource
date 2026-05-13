package alarm

import (
	"github.com/Compogo/compogo/component"
	"github.com/Compogo/compogo/container"
	"github.com/Compogo/compogo/flag"
)

var Component = &component.Component{
	Init: component.StepFunc(func(container container.Container) error {
		return container.Provide(NewConfig)
	}),
	BindFlags: component.BindFlags(func(flagSet flag.FlagSet, container container.Container) error {
		return container.Invoke(func(config *Config) {
			flagSet.Float32Var(&config.WarnStatePercent, WarnStatePercentFieldName, WarnStatePercentDefault, "")
			flagSet.Float32Var(&config.AlarmStatePercent, AlarmStatePercentFieldname, AlarmStatePercentDefault, "")
		})
	}),
	Configuration: component.StepFunc(func(container container.Container) error {
		return container.Invoke(Configuration)
	}),
}
