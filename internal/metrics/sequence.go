package metrics

import (
	"context"
	"time"

	"github.com/nokamoto/prometheus-mock-exporter/pkg/proto"
)

var afterFn = time.After

type sequence struct {
	sequence *proto.Sequence
}

func (s *sequence) run(ctx context.Context, mock *Mock) error {
	select {
	case <-ctx.Done():
		return nil
	case <-afterFn(s.sequence.GetInitialDelay().AsDuration()):
	}
	i := 0
	for {
		step := s.sequence.GetSteps()[i]
		mock.counters[step.GetId()].WithLabelValues(step.GetLabels()...).Add(step.GetValue())
		select {
		case <-ctx.Done():
			return nil
		case <-afterFn(1 * time.Second):
		}
		i = (i + 1) % len(s.sequence.GetSteps())
	}
}
