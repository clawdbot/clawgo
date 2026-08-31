package main

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestSleepContextCanceledReturnsCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := sleepContext(ctx, 30*time.Second)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("sleepContext(canceled, 30s) = %v, want context.Canceled", err)
	}
}

func TestSleepContextCancelDuringWait(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- sleepContext(ctx, 30*time.Second)
	}()
	cancel()

	var err error
	select {
	case err = <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("sleepContext did not return after cancel")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("sleepContext canceled mid-wait = %v, want context.Canceled", err)
	}
}

func TestSleepContextCompletesWhenContextStaysOpen(t *testing.T) {
	err := sleepContext(context.Background(), time.Millisecond)
	if err != nil {
		t.Fatalf("sleepContext(background, 1ms) = %v, want nil", err)
	}
}
