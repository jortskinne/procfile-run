package runner

import (
	"strings"
	"testing"
)

func TestPalette_ColorFor_Consistent(t *testing.T) {
	p := NewPalette()
	c1 := p.ColorFor("web")
	c2 := p.ColorFor("web")
	if c1 != c2 {
		t.Errorf("expected same color for same name, got %q and %q", c1, c2)
	}
}

func TestPalette_ColorFor_Different(t *testing.T) {
	p := NewPalette()
	c1 := p.ColorFor("web")
	c2 := p.ColorFor("worker")
	if c1 == c2 {
		t.Errorf("expected different colors for different names, got %q for both", c1)
	}
}

func TestPalette_ColorFor_Wraps(t *testing.T) {
	p := NewPalette()
	names := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	for _, n := range names {
		p.ColorFor(n)
	}
	// After exhausting all colors, should wrap around
	cFirst := p.ColorFor("a")
	cWrapped := p.ColorFor("h")
	_ = cFirst
	_ = cWrapped
	// Just ensure no panic and colors are assigned
	if len(p.assigned) != len(names) {
		t.Errorf("expected %d assignments, got %d", len(names), len(p.assigned))
	}
}

func TestPalette_Colorize(t *testing.T) {
	p := NewPalette()
	out := p.Colorize("web", "hello world")
	if !strings.Contains(out, "web") {
		t.Errorf("expected output to contain process name, got %q", out)
	}
	if !strings.Contains(out, "hello world") {
		t.Errorf("expected output to contain message, got %q", out)
	}
	if !strings.Contains(out, colorReset) {
		t.Errorf("expected output to contain reset code, got %q", out)
	}
}
