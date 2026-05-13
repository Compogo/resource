package alarm

import "github.com/Compogo/compogo/configurator"

const (
	WarnStatePercentFieldName  = "alarm.percent.warn"
	AlarmStatePercentFieldname = "alarm.percent.alarm"

	WarnStatePercentDefault  = 0.7
	AlarmStatePercentDefault = 0.85
)

type Config struct {
	WarnStatePercent  float32
	AlarmStatePercent float32
}

func NewConfig() *Config {
	return &Config{}
}

func Configuration(config *Config, configurator configurator.Configurator) *Config {
	if config.WarnStatePercent == 0 || config.WarnStatePercent == WarnStatePercentDefault {
		configurator.SetDefault(WarnStatePercentFieldName, WarnStatePercentDefault)
		config.WarnStatePercent = configurator.GetFloat32(WarnStatePercentFieldName)
	}

	if config.AlarmStatePercent == 0 || config.AlarmStatePercent == AlarmStatePercentDefault {
		configurator.SetDefault(AlarmStatePercentFieldname, AlarmStatePercentDefault)
		config.AlarmStatePercent = configurator.GetFloat32(AlarmStatePercentFieldname)
	}

	return config
}
