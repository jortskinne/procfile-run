package runner

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// Logger writes prefixed, colorized log lines for a named process.
type Logger struct {
	name    string
	palette *Palette
	out     io.Writer
}

// NewLogger creates a Logger that writes to out for the given process name.
func NewLogger(name string, palette *Palette, out io.Writer) *Logger {
	return &Logger{name: name, palette: palette, out: out}
}

// Write implements io.Writer so Logger can be used as a process stdout/stderr.
func (l *Logger) Write(p []byte) (n int, err error) {
	lines := strings.Split(strings.TrimRight(string(p), "\n"), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		ts := time.Now().Format("15:04:05")
		formatted := fmt.Sprintf("%s | %s\n", l.palette.Colorize(l.name, ts), line)
		_, err = fmt.Fprint(l.out, formatted)
		if err != nil {
			return 0, err
		}
	}
	return len(p), nil
}

// Log writes a plain message through the logger.
func (l *Logger) Log(msg string) {
	_, _ = l.Write([]byte(msg))
}
