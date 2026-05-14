package procfile

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Process represents a single entry in a Procfile.
type Process struct {
	Name    string
	Command string
}

// Parse reads a Procfile at the given path and returns a slice of Process entries.
// Lines starting with '#' or that are blank are ignored.
// Each valid line must be in the format: name: command
func Parse(path string) ([]Process, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("procfile: cannot open %q: %w", path, err)
	}
	defer f.Close()

	var processes []Process
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		idx := strings.Index(line, ":")
		if idx <= 0 {
			return nil, fmt.Errorf("procfile: invalid syntax on line %d: %q", lineNum, line)
		}

		name := strings.TrimSpace(line[:idx])
		command := strings.TrimSpace(line[idx+1:])

		if name == "" {
			return nil, fmt.Errorf("procfile: empty process name on line %d", lineNum)
		}
		if command == "" {
			return nil, fmt.Errorf("procfile: empty command for process %q on line %d", name, lineNum)
		}

		processes = append(processes, Process{Name: name, Command: command})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("procfile: error reading %q: %w", path, err)
	}

	if len(processes) == 0 {
		return nil, fmt.Errorf("procfile: no processes defined in %q", path)
	}

	return processes, nil
}
