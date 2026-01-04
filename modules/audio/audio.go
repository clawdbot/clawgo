package audio

import (
	"context"
	"time"
)

type Format struct {
	SampleRate int
	Channels   int
	Encoding   string
}

type Frame struct {
	Data      []byte
	Format    Format
	Timestamp time.Time
}

type Buffer struct {
	Data   []byte
	Format Format
}

type Capture interface {
	Name() string
	Start(ctx context.Context) (<-chan Frame, error)
	Close() error
}

type Playback interface {
	Name() string
	Play(ctx context.Context, in <-chan Frame) error
	Close() error
}
