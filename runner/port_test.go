package runner

import (
	"fmt"
	"net"
	"testing"
)

func TestExtractPort_Found(t *testing.T) {
	cases := []struct {
		cmd  string
		want int
	}{
		{"PORT=3000 node server.js", 3000},
		{"port=8080 python app.py", 8080},
		{"env PORT=5432 ./myapp", 5432},
	}
	for _, tc := range cases {
		port, ok := ExtractPort(tc.cmd)
		if !ok {
			t.Errorf("ExtractPort(%q): expected ok=true", tc.cmd)
		}
		if port != tc.want {
			t.Errorf("ExtractPort(%q): got %d, want %d", tc.cmd, port, tc.want)
		}
	}
}

func TestExtractPort_NotFound(t *testing.T) {
	cmds := []string{
		"node server.js",
		"./worker --verbose",
		"PORT= ./app",
	}
	for _, cmd := range cmds {
		_, ok := ExtractPort(cmd)
		if ok {
			t.Errorf("ExtractPort(%q): expected ok=false", cmd)
		}
	}
}

func TestIsPortInUse_Free(t *testing.T) {
	// Port 0 lets the OS pick a free port; we then check a high unlikely port.
	if IsPortInUse(19999) {
		t.Skip("port 19999 happens to be in use on this machine")
	}
}

func TestIsPortInUse_InUse(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal("could not open listener:", err)
	}
	defer ln.Close()
	port := ln.Addr().(*net.TCPAddr).Port
	if !IsPortInUse(port) {
		t.Errorf("expected port %d to be reported as in use", port)
	}
}

func TestCheckPorts_Conflict(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	port := ln.Addr().(*net.TCPAddr).Port

	processes := map[string]string{
		"web": fmt.Sprintf("PORT=%d node app.js", port),
		"worker": "./worker --queue default",
	}
	conflicts := CheckPorts(processes)
	if len(conflicts) != 1 {
		t.Errorf("expected 1 conflict, got %d: %v", len(conflicts), conflicts)
	}
}
