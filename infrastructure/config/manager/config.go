package manager

import (
	"time"

	"github.com/Compogo/compogo"
)

// DelayFieldName частота сбора статистики
const DelayFieldName = "manager.delay"

// DelayDefault 200ms по умолчанию
var DelayDefault = time.Second / 5

// Config содержит конфигурацию менеджера ресурсов.
type Config struct {
	Delay time.Duration
}

// NewConfig создаёт новую конфигурацию.
func NewConfig() *Config {
	return &Config{}
}

// Configuration загружает конфигурацию из Configurator.
func Configuration(config *Config, configurator compogo.Configurator) *Config {
	if config.Delay == 0 || config.Delay == DelayDefault {
		configurator.SetDefault(DelayFieldName, DelayDefault)
		config.Delay = configurator.GetDuration(DelayFieldName)
	}

	return config
}
