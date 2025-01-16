package config

import (
	"errors"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nokamoto/prometheus-mock-exporter/pkg/proto"
	"google.golang.org/protobuf/testing/protocmp"
)

func Test_LoadYamlConfig(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     *proto.Config
		wantErr  error
	}{
		{
			name:     "ok",
			filename: "testdata/config.yaml",
			want:     &proto.Config{},
		},
		{
			name:     "failed if file not found",
			filename: "testdata/notfound",
			wantErr:  os.ErrNotExist,
		},
		{
			name:     "failed if yaml is invalid",
			filename: "testdata/invalid",
			wantErr:  errUnmarshal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadYamlConfig(tt.filename)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("LoadYamlConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want, protocmp.Transform()); diff != "" {
				t.Errorf("LoadYamlConfig() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
