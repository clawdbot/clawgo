package stt

import "testing"

func TestParseBrabbleLineJSON(t *testing.T) {
	line := `{"event":"transcript","text":"hey razor"}`
	tr, ok := parseBrabbleLine(line)
	if !ok {
		t.Fatalf("expected ok")
	}
	if tr.Text != "hey razor" {
		t.Fatalf("expected text")
	}
	if !tr.Final {
		t.Fatalf("expected final")
	}
}

func TestParseBrabbleLinePayload(t *testing.T) {
	line := `{"event":"transcript","payload":{"transcript":"hello"}}`
	tr, ok := parseBrabbleLine(line)
	if !ok {
		t.Fatalf("expected ok")
	}
	if tr.Text != "hello" {
		t.Fatalf("expected payload text")
	}
}

func TestParseBrabbleLinePlain(t *testing.T) {
	line := "bring me on telegram"
	tr, ok := parseBrabbleLine(line)
	if !ok {
		t.Fatalf("expected ok")
	}
	if tr.Text != line {
		t.Fatalf("expected same line")
	}
}
