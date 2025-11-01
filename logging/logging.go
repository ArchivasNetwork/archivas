package logging

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

// ConfigureJSON configures the standard logger to emit JSONL records with a component tag.
func ConfigureJSON(component string) {
	log.SetFlags(0)
	log.SetOutput(&jsonWriter{
		component: component,
		out:       os.Stdout,
	})
}

type jsonWriter struct {
	component string
	out       io.Writer
	mu        sync.Mutex
}

func (w *jsonWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	message := strings.TrimRight(string(p), "\r\n")
	level, msg := parseLevel(message)

	entry := map[string]interface{}{
		"ts":        time.Now().UTC().Format(time.RFC3339Nano),
		"component": w.component,
		"level":     level,
		"msg":       msg,
	}

	encoded, err := json.Marshal(entry)
	if err != nil {
		return 0, err
	}

	encoded = append(encoded, '\n')
	if _, err := w.out.Write(encoded); err != nil {
		return 0, err
	}

	return len(p), nil
}

func parseLevel(message string) (string, string) {
	level := "info"
	msg := strings.TrimSpace(message)

	if strings.HasPrefix(msg, "[") {
		if idx := strings.Index(msg, "]"); idx > 0 && idx < 16 {
			token := strings.ToLower(strings.TrimSpace(msg[1:idx]))
			switch token {
			case "debug":
				level = "debug"
				msg = strings.TrimSpace(msg[idx+1:])
			case "info":
				level = "info"
				msg = strings.TrimSpace(msg[idx+1:])
			case "warn", "warning":
				level = "warn"
				msg = strings.TrimSpace(msg[idx+1:])
			case "error":
				level = "error"
				msg = strings.TrimSpace(msg[idx+1:])
			case "fatal":
				level = "fatal"
				msg = strings.TrimSpace(msg[idx+1:])
			default:
				// Keep default level when token is unrecognized
			}
		}
	}

	return level, msg
}
