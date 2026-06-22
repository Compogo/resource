package alarm

import (
	"github.com/Compogo/compogo"
)

const (
	// WarnStatePercentFieldName порог предупреждения
	WarnStatePercentFieldName = "alarm.percent.warn"

	// AlarmStatePercentFieldname порог тревоги
	AlarmStatePercentFieldname = "alarm.percent.alarm"
)

var (
	// WarnStatePercentDefault 70% — предупреждение
	WarnStatePercentDefault = float32(0.7)

	// AlarmStatePercentDefault 85% — тревога
	AlarmStatePercentDefault = float32(0.85)
)

// Config содержит конфигурацию алертов.
type Config struct {
	WarnStatePercent  float32
	AlarmStatePercent float32
}

// NewConfig создаёт новую конфигурацию.
func NewConfig() *Config {
	return &Config{}
}

// Configuration загружает конфигурацию из Configurator.
func Configuration(config *Config, configurator compogo.Configurator) *Config {
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
