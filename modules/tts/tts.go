package tts

import (
	"context"

	"github.com/clawdbot/clawgo/modules/audio"
)

type Request struct {
	Text   string
	Voice  string
	Rate   int
	Pitch  float32
	Engine string
}

type Engine interface {
	Name() string
	Synthesize(ctx context.Context, req Request) (audio.Buffer, error)
}
