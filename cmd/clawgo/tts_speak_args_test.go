package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
)

func TestSystemTTSEngineSpeakSeparatesUtteranceFromFlags(t *testing.T) {
	rec, out := buildArgRecorder(t)

	cases := []struct {
		name  string
		voice string
		rate  int
		text  string
	}{
		{name: "wav-write-flag", voice: "en-us", rate: 180, text: "-w/tmp/x"},
		{name: "short-unknown", voice: "en-us", rate: 180, text: "-foo"},
		{name: "markdown-list", voice: "en-us", rate: 180, text: "- item"},
		{name: "plain", voice: "en-us", rate: 180, text: "hello"},
		{name: "no-voice-rate", text: "-fpath"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := os.Remove(out); err != nil && !os.IsNotExist(err) {
				t.Fatalf("remove args file: %v", err)
			}
			engine := &systemTTSEngine{command: rec, voice: tc.voice, rate: tc.rate}
			if err := engine.Speak(tc.text); err != nil {
				t.Fatalf("speak: %v", err)
			}
			raw, err := os.ReadFile(out)
			if err != nil {
				t.Fatalf("read recorded args: %v", err)
			}
			var args []string
			if err := json.Unmarshal(raw, &args); err != nil {
				t.Fatalf("decode recorded args %q: %v", raw, err)
			}
			t.Logf("argv=%q", args)

			want := make([]string, 0, 6)
			if tc.voice != "" {
				want = append(want, "-v", tc.voice)
			}
			if tc.rate > 0 {
				want = append(want, "-s", strconv.Itoa(tc.rate))
			}
			want = append(want, "--", tc.text)
			if len(args) != len(want) {
				t.Fatalf("argv=%q want=%q", args, want)
			}
			for i := range want {
				if args[i] != want[i] {
					t.Fatalf("argv=%q want=%q", args, want)
				}
			}
		})
	}
}

func buildArgRecorder(t *testing.T) (bin, out string) {
	t.Helper()
	dir := t.TempDir()
	out = filepath.Join(dir, "args.json")
	src := filepath.Join(dir, "rec.go")
	srcText := `package main

import (
	"encoding/json"
	"os"
)

func main() {
	raw, err := json.Marshal(os.Args[1:])
	if err != nil {
		os.Exit(2)
	}
	if err := os.WriteFile(` + strconv.Quote(out) + `, raw, 0644); err != nil {
		os.Exit(3)
	}
}
`
	if err := os.WriteFile(src, []byte(srcText), 0644); err != nil {
		t.Fatalf("write recorder: %v", err)
	}
	bin = filepath.Join(dir, "rec")
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}
	cmd := exec.Command("go", "build", "-o", bin, src)
	cmd.Dir = dir
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build recorder: %v\n%s", err, output)
	}
	return bin, out
}
