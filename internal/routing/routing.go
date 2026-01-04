package routing

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

type Transport interface {
	SendVoiceTranscript(sessionKey, text string) error
	SendAgentRequest(sessionKey, text string, deliver bool, channel, to string) error
	SendProviderMessage(provider, to, message string) error
}

type Router interface {
	HandleTranscript(ctx context.Context, text string) (bool, error)
}

type Config struct {
	SessionKey       string
	AgentRequest     bool
	Deliver          bool
	DeliverChannel   string
	DeliverTo        string
	QuickActions     bool
	QuickPingMessage string
}

type Factory func(cfg Config, transport Transport, logf func(string, ...any)) (Router, error)

type registry struct {
	mu        sync.RWMutex
	factories map[string]Factory
}

var globalRegistry = &registry{factories: map[string]Factory{}}

func Register(name string, factory Factory) {
	name = strings.TrimSpace(name)
	if name == "" || factory == nil {
		return
	}
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()
	globalRegistry.factories[name] = factory
}

func New(name string, cfg Config, transport Transport, logf func(string, ...any)) (Router, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		name = "default"
	}
	globalRegistry.mu.RLock()
	factory, ok := globalRegistry.factories[name]
	globalRegistry.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("routing plugin not found: %s", name)
	}
	return factory(cfg, transport, logf)
}
