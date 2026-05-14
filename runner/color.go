package runner

import "fmt"

// ANSI color codes for process output
var colors = []string{
	"\033[36m", // cyan
	"\033[32m", // green
	"\033[33m", // yellow
	"\033[35m", // magenta
	"\033[34m", // blue
	"\033[31m", // red
	"\033[37m", // white
}

const colorReset = "\033[0m"

// Palette assigns a color to each process name by index.
type Palette struct {
	assigned map[string]string
	counter  int
}

// NewPalette creates a new Palette.
func NewPalette() *Palette {
	return &Palette{assigned: make(map[string]string)}
}

// ColorFor returns a consistent ANSI color string for the given process name.
func (p *Palette) ColorFor(name string) string {
	if c, ok := p.assigned[name]; ok {
		return c
	}
	c := colors[p.counter%len(colors)]
	p.assigned[name] = c
	p.counter++
	return c
}

// Colorize wraps text in the color assigned to name.
func (p *Palette) Colorize(name, text string) string {
	return fmt.Sprintf("%s%s%s %s", p.ColorFor(name), name, colorReset, text)
}
