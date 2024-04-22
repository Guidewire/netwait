package cli_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"gotest.tools/v3/icmd"
)

// CLI end-to-end tests
func TestCli(t *testing.T) {
	netwaitPath, err := buildCli()
	if err != nil {
		t.Fatal("failed to build app:", err)
	}

	// mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path == "/status/500" {
			writer.WriteHeader(500)
		}
	}))
	defer server.Close()
	serverUrl, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("failed to parse mock server URL: %s", err)
	}
	serverHostPort := serverUrl.Host

	tests := []struct {
		name     string
		cmd      icmd.Cmd
		expected icmd.Expected
	}{
		{
			name: "http",
			// netwait http://127.0.0.1:34287
			cmd: icmd.Command(netwaitPath, server.URL),
			expected: icmd.Expected{
				ExitCode: 0,
				// available: http://127.0.0.1:34287
				Out: "available: " + server.URL,
			},
		},
		{
			name: "http unavailable",
			// netwait http://127.0.0.1:34287/status/500 --timeout 5s
			cmd: icmd.Command(netwaitPath, server.URL+"/status/500", "--timeout", "5s"),
			expected: icmd.Expected{
				ExitCode: 1,
				// unavailable: http://127.0.0.1:34287/status/500
				Out: "unavailable: " + server.URL + "/status/500",
			},
		},
		{
			name: "tcp",
			// netwait 127.0.0.1:34287
			cmd: icmd.Command(netwaitPath, serverHostPort, "--timeout", "5s"),
			expected: icmd.Expected{
				ExitCode: 0,
				Out:      "available: " + serverHostPort,
			},
		},
		{
			name: "dns",
			// netwait localhost
			cmd: icmd.Command(netwaitPath, "localhost"),
			expected: icmd.Expected{
				ExitCode: 0,
				Out:      "available: localhost",
			},
		},
		{
			name: "multiple resources",
			// netwait http://127.0.0.1:34287 127.0.0.1:34287
			cmd: icmd.Command(netwaitPath, server.URL, serverHostPort),
			expected: icmd.Expected{
				ExitCode: 0,
				// available: http://127.0.0.1:34287
				// available: 127.0.0.1:34287
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(os.Getwd())
			result := icmd.RunCmd(tt.cmd)
			result.Assert(t, tt.expected)
		})
	}
}

// Build app, return path to executable
func buildCli() (string, error) {
	cmd := exec.Command("go", "build", "-o", "build/netwait")
	cmd.Dir = ".."
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("build command failed: %w", err)
	}
	path, _ := filepath.Abs("../build/netwait")
	_, err = os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("file does not exist: %w", err)
	}
	return path, nil
}
