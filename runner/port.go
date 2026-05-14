package runner

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
)

var portRegex = regexp.MustCompile(`(?i)(?:PORT|port)\s*=\s*(\d+)`)

// ExtractPort attempts to find a PORT=XXXX assignment in a command string.
func ExtractPort(command string) (int, bool) {
	matches := portRegex.FindStringSubmatch(command)
	if len(matches) < 2 {
		return 0, false
	}
	port, err := strconv.Atoi(matches[1])
	if err != nil || port < 1 || port > 65535 {
		return 0, false
	}
	return port, true
}

// IsPortInUse checks whether the given TCP port is already bound on localhost.
func IsPortInUse(port int) bool {
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return true
	}
	ln.Close()
	return false
}

// CheckPorts inspects a map of process name → command and returns a list of
// conflict descriptions for any ports that are already in use.
func CheckPorts(processes map[string]string) []string {
	var conflicts []string
	for name, cmd := range processes {
		port, ok := ExtractPort(cmd)
		if !ok {
			continue
		}
		if IsPortInUse(port) {
			conflicts = append(conflicts, fmt.Sprintf("%s: port %d is already in use", name, port))
		}
	}
	return conflicts
}
