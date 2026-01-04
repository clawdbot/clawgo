package stt

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"

	"github.com/clawdbot/clawgo/modules/audio"
)

type BrabbleConfig struct {
	Command string
	Args    []string
}

type BrabbleEngine struct {
	cfg  BrabbleConfig
	logf func(string, ...any)
}

func NewBrabbleEngine(cfg BrabbleConfig, logf func(string, ...any)) Engine {
	return &BrabbleEngine{cfg: cfg, logf: logf}
}

func (e *BrabbleEngine) Name() string { return "brabble" }

func (e *BrabbleEngine) Transcribe(ctx context.Context, _ <-chan audio.Frame, _ Options) (<-chan Transcript, error) {
	cmdPath := strings.TrimSpace(e.cfg.Command)
	if cmdPath == "" {
		cmdPath = "brabble"
	}
	cmd := exec.CommandContext(ctx, cmdPath, e.cfg.Args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	out := make(chan Transcript, 32)
	go func() {
		defer close(out)
		e.readLines(ctx, stdout, out)
	}()
	go e.logLines(ctx, stderr)
	go func() {
		_ = cmd.Wait()
	}()
	return out, nil
}

func (e *BrabbleEngine) readLines(ctx context.Context, r io.Reader, out chan<- Transcript) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 8*1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		tr, ok := parseBrabbleLine(line)
		if !ok {
			continue
		}
		select {
		case out <- tr:
		case <-ctx.Done():
			return
		}
	}
	if err := scanner.Err(); err != nil && !errors.Is(err, io.EOF) && e.logf != nil {
		e.logf("brabble read error: %v", err)
	}
}

func (e *BrabbleEngine) logLines(ctx context.Context, r io.Reader) {
	if e.logf == nil {
		return
	}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return
		default:
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		e.logf("brabble: %s", line)
	}
}

type brabbleEvent struct {
	Type       string          `json:"type"`
	Event      string          `json:"event"`
	Text       string          `json:"text"`
	Transcript string          `json:"transcript"`
	Utterance  string          `json:"utterance"`
	Final      *bool           `json:"final"`
	Payload    json.RawMessage `json:"payload"`
}

type brabblePayload struct {
	Text       string `json:"text"`
	Transcript string `json:"transcript"`
	Utterance  string `json:"utterance"`
}

func parseBrabbleLine(line string) (Transcript, bool) {
	if strings.HasPrefix(line, "{") {
		var evt brabbleEvent
		if err := json.Unmarshal([]byte(line), &evt); err == nil {
			text := pickBrabbleText(evt.Text, evt.Transcript, evt.Utterance)
			if text == "" && len(evt.Payload) > 0 {
				var payload brabblePayload
				if err := json.Unmarshal(evt.Payload, &payload); err == nil {
					text = pickBrabbleText(payload.Text, payload.Transcript, payload.Utterance)
				}
			}
			if text == "" {
				return Transcript{}, false
			}
			final := true
			if evt.Final != nil {
				final = *evt.Final
			}
			if strings.Contains(strings.ToLower(evt.Type), "partial") || strings.Contains(strings.ToLower(evt.Event), "partial") {
				final = false
			}
			return Transcript{Text: text, Final: final, Timestamp: time.Now(), Source: "brabble"}, true
		}
	}
	return Transcript{Text: strings.TrimSpace(line), Final: true, Timestamp: time.Now(), Source: "brabble"}, true
}

func pickBrabbleText(parts ...string) string {
	for _, part := range parts {
		if strings.TrimSpace(part) != "" {
			return strings.TrimSpace(part)
		}
	}
	return ""
}

func (e *BrabbleEngine) String() string {
	cmd := strings.TrimSpace(e.cfg.Command)
	if cmd == "" {
		cmd = "brabble"
	}
	if len(e.cfg.Args) == 0 {
		return cmd
	}
	return fmt.Sprintf("%s %s", cmd, strings.Join(e.cfg.Args, " "))
}
