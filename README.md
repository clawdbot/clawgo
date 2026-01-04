# Clawgo (Go Node)

Minimal headless node client for Raspberry Pi / Linux. Connects to the gateway bridge, handles pairing, streams `voice.transcript` events (stdin/FIFO), subscribes to chat, and can speak responses via local TTS.

## Build

```bash
cd clawgo
go build ./cmd/clawgo
```

Cross-compile for Pi:

```bash
GOOS=linux GOARCH=arm64 go build -o /tmp/clawgo-linux-arm64 ./cmd/clawgo
```

## Key flags

| Flag | Description |
| --- | --- |
| `-session-key` | Session for outgoing `voice.transcript` events (default `main`). |
| `-chat-session-key` | Session to subscribe for chat replies (default mirrors `-session-key`). |
| `-chat-subscribe` | Enable chat stream+TTS (default `true`). |
| `-tts-engine` | `system`, `piper`, `elevenlabs`, or `none` (system = `espeak-ng`). |
| `-tts-system-voice` | espeak voice id (default `en-us`). |
| `-tts-system-rate` | Speech rate (wpm). |
| `-mdns-service` | Bonjour service type (default `_clawdis-node._tcp`). |
| `-stdin` | Read transcripts from stdin (pipe/FIFO). |
| `-stdin-file` | Read transcripts from a FIFO/file instead of stdin. |
| `-agent-request` | Send transcripts as `agent.request` (uses agent + deliver). |
| `-deliver` | Deliver agent responses to a provider (requires channel + to). |
| `-deliver-channel` | Delivery provider (telegram/whatsapp/signal/imessage). |
| `-deliver-to` | Delivery destination id. |
| `-quick-actions` | Enable built-in quick actions (default true). |
| `-ping-message` | Message used for telegram ping quick action. |

## Pair

```bash
./clawgo pair \
  -bridge 100.88.46.29:18790 \
  -display-name "Razor Pi"
```

Approve via `clawdis nodes approve <requestId>`.

## Run (FIFO + TTS example)

```bash
mkfifo /tmp/voice.fifo
# in one terminal
tail -f /tmp/voice.fifo | ./clawgo run \ 
  -bridge 100.88.46.29:18790 \
  -stdin \
  -chat-subscribe \
  -tts-engine system
# elsewhere
printf hey computer turn on the lightsn > /tmp/voice.fifo
```

Each line on the FIFO becomes a `voice.transcript`; chat responses from the `main` session are spoken via `espeak-ng`.

## systemd example

See `docs/linux-node.md` for the end-to-end Pi setup. TL;DR:

1. Install the binary as `/home/pi/clawgo`.
2. Create a wrapper script that keeps a FIFO (`/home/pi/.cache/clawdis/voice.fifo`) open and pipes it into `clawgo run -stdin`.
3. Create `/etc/systemd/system/clawgo.service` pointing to that wrapper.

## mDNS advertising

The node advertises `_clawdis-node._tcp` by default.

```bash
dns-sd -B _clawdis-node._tcp local.
```

Override to `_clawdis-bridge._tcp` if you intentionally want it to show up as a gateway beacon:

```bash
./clawgo run -mdns-service _clawdis-bridge._tcp
```

## Notes
- Node state (`nodeId` + token) lives in `~/.clawdis/clawgo.json`.
- Caps default to `voiceWake`; override via `-caps` if you expose more commands.
- Set `bridge.bind: "tailnet"` on the gateway to restrict the bridge to Tailscale.
