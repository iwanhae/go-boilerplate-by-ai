package logging

import (
	"strings"
	"sync"
)

// Buffer collects log lines for later retrieval.
type Buffer struct {
	mu    sync.Mutex
	lines []string
}

// Write implements io.Writer, storing each line.
func (b *Buffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.lines = append(b.lines, strings.TrimSpace(string(p)))
	return len(p), nil
}

// Lines returns the buffered log lines.
func (b *Buffer) Lines() []string {
	b.mu.Lock()
	defer b.mu.Unlock()
	out := make([]string, len(b.lines))
	copy(out, b.lines)
	return out
}
