package metric

import (
	"context"

	"github.com/Compogo/compogo"
	"github.com/Compogo/resource/domain"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	CpuLabel    = "cpu"
	MemoryLabel = "memory"
)

type Metric struct {
	limit *prometheus.GaugeVec
	usage *prometheus.GaugeVec
}

func NewMetric(compogoConfig *compogo.Config) *Metric {
	return &Metric{
		limit: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name:        compogo.MetricNamePrefix + "resource_limit",
			Help:        "application resource limits",
			ConstLabels: prometheus.Labels{compogo.MetricAppNameFieldName: compogoConfig.Name},
		}, []string{CpuLabel, MemoryLabel}),
		usage: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name:        compogo.MetricNamePrefix + "resource_usage",
			Help:        "resources used by the application",
			ConstLabels: prometheus.Labels{compogo.MetricAppNameFieldName: compogoConfig.Name},
		}, []string{CpuLabel, MemoryLabel}),
	}
}

func (metric *Metric) OnChangeResource(_ context.Context, resource *domain.Resource) {
	switch resource.Type {
	case domain.CPU:
		metric.limit.WithLabelValues(CpuLabel).Set(float64(resource.Limit))
		metric.usage.WithLabelValues(CpuLabel).Set(float64(resource.Usage))
	case domain.Memory:
		metric.limit.WithLabelValues(MemoryLabel).Set(float64(resource.Limit))
		metric.usage.WithLabelValues(MemoryLabel).Set(float64(resource.Usage))
	}
}
