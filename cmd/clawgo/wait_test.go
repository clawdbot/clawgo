package main

import (
	"context"
	"errors"
	"testing"
	"time"
)

func testBridgeClient() *BridgeClient {
	return &BridgeClient{
		logf:   func(string, ...any) {},
		errs:   make(chan error),
		frames: make(chan map[string]any),
	}
}

func TestWaitForPairCancelUnblocks(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	c := testBridgeClient()
	done := make(chan error, 1)
	go func() {
		_, err := waitForPair(ctx, c)
		done <- err
	}()
	cancel()

	var err error
	select {
	case err = <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("waitForPair did not return after cancel")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("waitForPair canceled = %v, want context.Canceled", err)
	}
}

func TestWaitForHelloCancelUnblocks(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	c := testBridgeClient()
	done := make(chan error, 1)
	go func() {
		done <- waitForHello(ctx, c)
	}()
	cancel()

	var err error
	select {
	case err = <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("waitForHello did not return after cancel")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("waitForHello canceled = %v, want context.Canceled", err)
	}
}

func TestWaitForPairReturnsTokenOnPairOK(t *testing.T) {
	c := testBridgeClient()
	c.frames = make(chan map[string]any, 1)
	c.frames <- map[string]any{"type": "pair-ok", "token": "tok-1"}

	token, err := waitForPair(context.Background(), c)
	if err != nil {
		t.Fatalf("waitForPair pair-ok: %v", err)
	}
	if token != "tok-1" {
		t.Fatalf("waitForPair token = %q, want tok-1", token)
	}
}

func TestWaitForHelloReturnsOnHelloOK(t *testing.T) {
	c := testBridgeClient()
	c.frames = make(chan map[string]any, 1)
	c.frames <- map[string]any{"type": "hello-ok", "serverName": "gw"}

	if err := waitForHello(context.Background(), c); err != nil {
		t.Fatalf("waitForHello hello-ok: %v", err)
	}
}
