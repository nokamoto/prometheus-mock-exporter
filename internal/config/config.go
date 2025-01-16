package config

import (
	"fmt"
	"os"

	"buf.build/go/protoyaml"
	"github.com/nokamoto/prometheus-mock-exporter/pkg/proto"
)

var errUnmarshal = fmt.Errorf("failed to unmarshal")

// LoadYamlConfig loads a proto.Config from a yaml file.
func LoadYamlConfig(filename string) (*proto.Config, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", filename, err)
	}

	var config proto.Config
	if err := protoyaml.Unmarshal(bytes, &config); err != nil {
		return nil, fmt.Errorf("%w %s: %w", errUnmarshal, filename, err)
	}
	return &config, nil
}
