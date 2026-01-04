package stt

import (
	"context"
	"strings"
	"time"

	"github.com/clawdbot/clawgo/modules/audio"
)

type LineEngine struct{}

func NewLineEngine() Engine { return &LineEngine{} }

func (e *LineEngine) Name() string { return "line" }

func (e *LineEngine) Transcribe(ctx context.Context, in <-chan audio.Frame, _ Options) (<-chan Transcript, error) {
	out := make(chan Transcript, 32)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case frame, ok := <-in:
				if !ok {
					return
				}
				text := strings.TrimSpace(string(frame.Data))
				if text == "" {
					continue
				}
				ts := frame.Timestamp
				if ts.IsZero() {
					ts = time.Now()
				}
				tr := Transcript{Text: text, Final: true, Timestamp: ts, Source: "line"}
				select {
				case out <- tr:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return out, nil
}
