package stt

import (
	"context"
	"time"

	"github.com/clawdbot/clawgo/modules/audio"
)

type Transcript struct {
	Text       string
	Final      bool
	Confidence float32
	Timestamp  time.Time
	Source     string
}

type Options struct {
	Language string
	Prompt   string
	Model    string
}

type Engine interface {
	Name() string
	Transcribe(ctx context.Context, in <-chan audio.Frame, opts Options) (<-chan Transcript, error)
}
