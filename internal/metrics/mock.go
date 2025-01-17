package metrics

import (
	"errors"
	"fmt"

	"github.com/nokamoto/prometheus-mock-exporter/pkg/proto"
	"github.com/prometheus/client_golang/prometheus"
)

type Mock struct {
	counters map[string]*prometheus.CounterVec
}

var errDuplicateID = errors.New("duplicate id")

// New creates a new Mock from the given configuration.
//
// All metrics are automatically registered to the default registry.
func New(config *proto.Config) (*Mock, error) {
	mock := &Mock{
		counters: make(map[string]*prometheus.CounterVec),
	}

	for _, counter := range config.GetCounters() {
		key := counter.GetId()
		if _, ok := mock.counters[key]; ok {
			return nil, fmt.Errorf("%w: counter %s", errDuplicateID, key)
		}

		mock.counters[key] = prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: counter.GetNamespace(),
			Name:      counter.GetName(),
			Help:      counter.GetHelp(),
		}, counter.GetLabels())
	}

	return mock, nil
}

// MustRegister registers all metrics to the given registry.
func (m *Mock) MustRegister(registry *prometheus.Registry) {
	for _, counter := range m.counters {
		registry.MustRegister(counter)
	}
}
