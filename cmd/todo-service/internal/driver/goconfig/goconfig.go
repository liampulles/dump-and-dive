package goconfig

import (
	"fmt"

	"github.com/liampulles/go-config"

	"github.com/liampulles/dump-and-dive/cmd/todo-service/internal/adapter"
)

// Provider provides config properties via the go-config module
type Provider struct {
	source config.Source
}

var _ adapter.ConfigProvider = &Provider{}

// NewProvider is a constructor
func NewProvider(source config.Source) *Provider {
	return &Provider{
		source: source,
	}
}

type configSnapshot struct {
	port int
}

var _ adapter.ConfigSnapshot = &configSnapshot{}

// Port implements adapter.ConfigSnapshot
func (cs *configSnapshot) Port() int {
	return cs.port
}

// GetAdapterConfig implements adapter.ConfigProvider
func (gcp *Provider) GetAdapterConfig() (adapter.ConfigSnapshot, error) {
	return gcp.get()
}

func (gcp *Provider) get() (*configSnapshot, error) {
	typedSource := config.NewTypedSource(gcp.source)
	result := &configSnapshot{
		port: 9090,
	}
	if err := config.LoadProperties(typedSource,
		config.IntProp("PORT", &result.port, false),
	); err != nil {
		return nil, fmt.Errorf("could not fetch config: %w", err)
	}
	return result, nil
}
