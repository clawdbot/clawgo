package main

import (
	"testing"
	"time"
)

type nopTTSEngine struct{}

func (nopTTSEngine) Speak(string) error { return nil }

func TestReplaceTTSQueueStopsPreviousLoop(t *testing.T) {
	first := newTTSQueue(nopTTSEngine{}, func(string, ...any) {})
	if first == nil {
		t.Fatal("newTTSQueue returned nil")
	}
	second := replaceTTSQueue(first, nopTTSEngine{}, func(string, ...any) {})
	if second == nil {
		t.Fatal("replaceTTSQueue returned nil")
	}
	t.Cleanup(second.Stop)

	select {
	case <-first.done:
	case <-time.After(2 * time.Second):
		t.Fatal("previous TTS queue still ranging after reconnect")
	}
}
