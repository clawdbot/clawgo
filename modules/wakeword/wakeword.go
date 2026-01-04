package wakeword

import (
	"context"
	"time"

	"github.com/clawdbot/clawgo/modules/audio"
)

type Event struct {
	Keyword    string
	Confidence float32
	Timestamp  time.Time
}

type Detector interface {
	Name() string
	Detect(ctx context.Context, in <-chan audio.Frame) (<-chan Event, error)
}
