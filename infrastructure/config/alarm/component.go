package alarm

import (
	"github.com/Compogo/compogo"
	"github.com/Compogo/compogo/flag"
)

// Component — компонент конфигурации алертов.
// Регистрирует конфигурацию в DI-контейнере.
var Component = compogo.Component{
	Name: "resource.alarm.config",
	Init: compogo.StepFunc(func(container compogo.Container) error {
		return container.Provide(NewConfig)
	}),
	BindFlags: compogo.BindFlags(func(flagSet flag.FlagSet, container compogo.Container) error {
		return container.Invoke(func(config *Config) {
			flagSet.Float32Var(&config.WarnStatePercent, WarnStatePercentFieldName, WarnStatePercentDefault, "")
			flagSet.Float32Var(&config.AlarmStatePercent, AlarmStatePercentFieldname, AlarmStatePercentDefault, "")
		})
	}),
	Configuration: compogo.StepFunc(func(container compogo.Container) error {
		return container.Invoke(Configuration)
	}),
}
