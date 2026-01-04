package defaultpolicy

import (
	"context"
	"strings"
	"unicode"

	"github.com/clawdbot/clawgo/internal/routing"
)

type Policy struct {
	cfg       routing.Config
	transport routing.Transport
	logf      func(string, ...any)
}

func init() {
	routing.Register("default", New)
}

func New(cfg routing.Config, transport routing.Transport, logf func(string, ...any)) (routing.Router, error) {
	if logf == nil {
		logf = func(string, ...any) {}
	}
	return &Policy{cfg: cfg, transport: transport, logf: logf}, nil
}

func (p *Policy) HandleTranscript(_ context.Context, text string) (bool, error) {
	if strings.TrimSpace(text) == "" {
		return false, nil
	}
	if p.cfg.QuickActions {
		if handled, err := p.handleQuickActions(text); handled {
			return true, err
		}
	}
	if p.cfg.AgentRequest {
		return true, p.transport.SendAgentRequest(p.cfg.SessionKey, text, p.cfg.Deliver, p.cfg.DeliverChannel, p.cfg.DeliverTo)
	}
	return true, p.transport.SendVoiceTranscript(p.cfg.SessionKey, text)
}

func (p *Policy) handleQuickActions(text string) (bool, error) {
	if isTelegramPing(text) {
		if p.cfg.DeliverChannel != "telegram" || p.cfg.DeliverTo == "" {
			p.logf("quick action skipped: telegram delivery not configured")
			return true, nil
		}
		message := strings.TrimSpace(p.cfg.QuickPingMessage)
		if message == "" {
			message = "Ping."
		}
		err := p.transport.SendProviderMessage("telegram", p.cfg.DeliverTo, message)
		if err != nil {
			p.logf("quick action send failed: %v", err)
		} else {
			p.logf("quick action: telegram ping sent")
		}
		return true, err
	}
	return false, nil
}

func isTelegramPing(text string) bool {
	normalized := normalizeCommand(text)
	return strings.Contains(normalized, "telegram") && strings.Contains(normalized, "ping")
}

func normalizeCommand(text string) string {
	lowered := strings.ToLower(text)
	var b strings.Builder
	for _, r := range lowered {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			b.WriteRune(r)
			continue
		}
		b.WriteRune(32)
	}
	return strings.Join(strings.Fields(b.String()), " ")
}
