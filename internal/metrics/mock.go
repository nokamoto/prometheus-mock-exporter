package metrics

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/nokamoto/prometheus-mock-exporter/pkg/proto"
	"github.com/prometheus/client_golang/prometheus"
)

type Mock struct {
	counters  map[string]*prometheus.CounterVec
	sequences []*sequence
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

	for _, s := range config.GetSequences() {
		mock.sequences = append(mock.sequences, &sequence{
			sequence: s,
		})
	}

	return mock, nil
}

// MustRegister registers all metrics to the given registry.
func (m *Mock) MustRegister(registry *prometheus.Registry) {
	for _, counter := range m.counters {
		registry.MustRegister(counter)
	}
}

// Run starts all sequences concurrently until the given context is canceled.
func (m *Mock) Run(ctx context.Context) {
	var wg sync.WaitGroup
	for _, s := range m.sequences {
		wg.Add(1)
		go func ()  {
			defer wg.Done()
			s.run(ctx, m)
		}()
	}
	wg.Wait()
}
