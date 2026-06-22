# Compogo Resource

[![Go Reference](https://pkg.go.dev/badge/github.com/Compogo/resource.svg)](https://pkg.go.dev/github.com/Compogo/resource)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Мониторинг и управление системными ресурсами для фреймворка [Compogo](https://github.com/Compogo/compogo).

Предоставляет:

* Мониторинг CPU и Memory через cgroup
* Автоматическое применение лимитов к рантайму Go
* Систему алертов (Warning/Alarm)
* Метрики Prometheus
* Событийную модель для интеграции

## Установка

```shell
go get github.com/Compogo/resource
```

## Быстрый старт

```go
package main

import (
    "github.com/Compogo/compogo"
    "github.com/Compogo/resource/application/alarm"
    "github.com/Compogo/resource/application/manager"
    "github.com/Compogo/resource/application/metric"
    "github.com/Compogo/resource/application/quota"
)

func main() {
    app := compogo.NewApp("myapp",
        // Менеджер ресурсов
        compogo.WithComponents(&manager.Component),
        
        // Применение квот к рантайму
        compogo.WithComponents(&quota.Component),
        
        // Система алертов
        compogo.WithComponents(&alarm.Component),
        
        // Метрики Prometheus
        compogo.WithComponents(&metric.Component),
    )

    if err := app.Serve(); err != nil {
        panic(err)
    }
}
```

## Конфигурация

### Алерты

```shell
# Порог предупреждения (70% по умолчанию)
--alarm.percent.warn=0.7

# Порог тревоги (85% по умолчанию)
--alarm.percent.alarm=0.85
```

### Менеджер ресурсов

```shell
# Частота сбора статистики (200ms по умолчанию)
--manager.delay=200ms
```

## Использование

### Мониторинг ресурсов

```go
var m *manager.Manager
container.Invoke(func(manager *manager.Manager) { m = manager })

// Текущий CPU
cpu := m.CPU()
fmt.Printf("CPU: limit=%d, usage=%d\n", cpu.Limit, cpu.Usage)

// Текущая память
memory := m.Memory()
fmt.Printf("Memory: limit=%d, usage=%d\n", memory.Limit, memory.Usage)
```

### Подписка на изменения

```go
type MyComponent struct {
    logger Logger
}

func (c *MyComponent) OnChangeResource(ctx context.Context, resource *domain.Resource) {
    c.logger.Infof("Resource changed: type=%s, limit=%d, usage=%d", 
        resource.Type, resource.Limit, resource.Usage)
}

// Подписка
m.OnChangeResource.Subscribe(c.OnChangeResource)
```

### Алерты

```go
type AlertHandler struct {
    logger Logger
}

func (h *AlertHandler) OnAlarm(ctx context.Context, state domain.State) {
    switch state {
    case domain.Warning:
        h.logger.Warn("Resource usage exceeded 70%")
    case domain.Alarm:
        h.logger.Error("Resource usage exceeded 85%")
    }
}

// Подписка на алерты
alarm.OnAlarm.Subscribe(h.OnAlarm)

// Текущее состояние
state := alarm.State()
```

## Автоматическое применение лимитов

```go
// Компонент quota автоматически применяет лимиты к рантайму:
// - GOMAXPROCS = CPU limit / 100000
// - SetMemoryLimit = Memory limit
```

## Метрики Prometheus

```text
# Лимиты ресурсов
compogo_resource_limit{app="myapp", cpu="cpu"} 200000
compogo_resource_limit{app="myapp", memory="memory"} 1073741824

# Использование ресурсов
compogo_resource_usage{app="myapp", cpu="cpu"} 45000
compogo_resource_usage{app="myapp", memory="memory"} 536870912
```

Зависимости

* [Compogo](https://github.com/Compogo/compogo) — основной фреймворк
* [Compogo Repeater](https://github.com/Compogo/repeater) — периодический сбор статистики
* [Compogo Types](https://github.com/Compogo/types) — утилиты (emitter)
* [Prometheus](https://github.com/prometheus/client_golang) — метрики

## Лицензия

```text
MIT License

Copyright (c) 2026 Compogo

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

```
