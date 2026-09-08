package main

import (
	"testing"

	"github.com/clawdbot/clawgo/modules/stt"
)

func TestShouldForwardTranscript(t *testing.T) {
	cases := []struct {
		name string
		tr   stt.Transcript
		want bool
	}{
		{name: "non-final", tr: stt.Transcript{Text: "hey razor", Final: false}, want: false},
		{name: "final", tr: stt.Transcript{Text: "hey razor", Final: true}, want: true},
		{name: "empty-final", tr: stt.Transcript{Text: "  ", Final: true}, want: false},
		{name: "empty-non-final", tr: stt.Transcript{Text: "", Final: false}, want: false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := shouldForwardTranscript(tc.tr)
			t.Logf("shouldForwardTranscript(Text=%q Final=%v) = %v", tc.tr.Text, tc.tr.Final, got)
			if got != tc.want {
				t.Fatalf("shouldForwardTranscript(Text=%q Final=%v)=%v want %v", tc.tr.Text, tc.tr.Final, got, tc.want)
			}
		})
	}
}
