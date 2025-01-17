package metrics

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/nokamoto/prometheus-mock-exporter/pkg/proto"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		config  *proto.Config
		wantErr error
	}{
		{
			name: "ok",
			config: &proto.Config{
				Counters: []*proto.Counter{
					{
						Id: "c1",
					},
					{
						Id: "c2",
					},
				},
			},
		},
		{
			name: "failed if counter id is duplicated",
			config: &proto.Config{
				Counters: []*proto.Counter{
					{
						Id: "c1",
					},
					{
						Id: "c1",
					},
				},
			},
			wantErr: errDuplicateID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := New(tt.config)
			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("New() error = %v, wantErr %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestMock_Run(t *testing.T) {
	tests := []struct {
		name   string
		config *proto.Config
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := New(tt.config)
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}
			mock.Run(context.TODO())
		})
	}
}
