package runner

import (
	"bufio"
	"os"
	"strings"
)

// LoadEnvFile reads a .env file and returns a map of key=value pairs.
// Lines starting with '#' and empty lines are ignored.
// Existing environment variables are not overwritten.
func LoadEnvFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]string{}, nil
		}
		return nil, err
	}
	defer f.Close()

	env := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := parseLine(line)
		if !ok {
			continue
		}
		env[key] = value
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return env, nil
}

// parseLine splits a KEY=VALUE line into its parts.
// Inline comments (# ...) after the value are stripped.
// Surrounding quotes on the value are removed.
func parseLine(line string) (string, string, bool) {
	idx := strings.IndexByte(line, '=')
	if idx <= 0 {
		return "", "", false
	}
	key := strings.TrimSpace(line[:idx])
	val := strings.TrimSpace(line[idx+1:])

	// Strip inline comment
	if ci := strings.Index(val, " #"); ci >= 0 {
		val = strings.TrimSpace(val[:ci])
	}

	// Strip surrounding quotes
	if len(val) >= 2 {
		if (val[0] == '"' && val[len(val)-1] == '"') ||
			(val[0] == '\'' && val[len(val)-1] == '\'') {
			val = val[1 : len(val)-1]
		}
	}

	if key == "" {
		return "", "", false
	}
	return key, val, true
}

// EnvMapToSlice converts a map of env vars to a []string slice
// in the form KEY=VALUE, suitable for exec.Cmd.Env.
func EnvMapToSlice(env map[string]string) []string {
	slice := make([]string, 0, len(env))
	for k, v := range env {
		slice = append(slice, k+"="+v)
	}
	return slice
}
