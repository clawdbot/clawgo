package main

import (
	"os/exec"
	"testing"
	"time"
)

func TestSystemTTSEngineSpeakHungProcessIsBounded(t *testing.T) {
	sleep, err := exec.LookPath("sleep")
	if err != nil {
		t.Skip("sleep not found")
	}

	engine := &systemTTSEngine{
		command: sleep,
		timeout: 200 * time.Millisecond,
	}
	if err := engine.Speak("2"); err == nil {
		t.Fatal("expected hung speak to return an error")
	}
}
