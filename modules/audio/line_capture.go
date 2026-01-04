package audio

import (
	"bufio"
	"context"
	"io"
	"os"
	"strings"
	"time"
)

type LineCapture struct {
	name   string
	path   string
	read   io.Reader
	logf   func(string, ...any)
	format Format
}

func NewLineCapture(name string, reader io.Reader, logf func(string, ...any)) *LineCapture {
	return &LineCapture{
		name:   name,
		read:   reader,
		logf:   logf,
		format: Format{Encoding: "text/line"},
	}
}

func NewLineCaptureFromPath(path string, logf func(string, ...any)) *LineCapture {
	name := strings.TrimSpace(path)
	if name == "" {
		name = "fifo"
	}
	return &LineCapture{
		name:   name,
		path:   path,
		logf:   logf,
		format: Format{Encoding: "text/line"},
	}
}

func (c *LineCapture) Name() string {
	if strings.TrimSpace(c.name) == "" {
		return "line"
	}
	return c.name
}

func (c *LineCapture) Start(ctx context.Context) (<-chan Frame, error) {
	out := make(chan Frame, 32)
	if c.path != "" {
		go c.readPathLoop(ctx, out)
	} else {
		go c.readReader(ctx, c.read, out)
	}
	return out, nil
}

func (c *LineCapture) Close() error { return nil }

func (c *LineCapture) readReader(ctx context.Context, r io.Reader, out chan<- Frame) {
	defer close(out)
	if r == nil {
		return
	}
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 8*1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		frame := Frame{Data: []byte(line), Format: c.format, Timestamp: time.Now()}
		select {
		case out <- frame:
		case <-ctx.Done():
			return
		}
	}
	if err := scanner.Err(); err != nil && c.logf != nil {
		c.logf("capture read error: %v", err)
	}
}

func (c *LineCapture) readPathLoop(ctx context.Context, out chan<- Frame) {
	defer close(out)
	path := strings.TrimSpace(c.path)
	if path == "" {
		return
	}
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		f, err := os.Open(path)
		if err != nil {
			if c.logf != nil {
				c.logf("capture open failed: %v", err)
			}
			time.Sleep(500 * time.Millisecond)
			continue
		}
		c.readReader(ctx, f, out)
		_ = f.Close()
		time.Sleep(200 * time.Millisecond)
	}
}
