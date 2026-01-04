package tts

type Request struct {
	Text   string
	Voice  string
	Rate   int
	Engine string
}
