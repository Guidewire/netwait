package cli_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"

	"gotest.tools/v3/icmd"
)

func TestCli(t *testing.T) {
	cmd := exec.Command("go", "build", "-o", "cli_test/bin/netwait")
	cmd.Dir = ".."
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}

	netwaitPath := "bin/netwait"
	_, err = os.Stat(netwaitPath)
	if err != nil {
		t.Fatal("netwait binary is required", err)
	}

	// mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	}))
	defer server.Close()

	// mock HTTP server, returns error
	serverFail := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(500)
	}))
	defer serverFail.Close()

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
			// netwait http://127.0.0.1:34287 --timeout 5s
			cmd: icmd.Command(netwaitPath, serverFail.URL, "--timeout", "5s"),
			expected: icmd.Expected{
				ExitCode: 1,
				// unavailable: http://127.0.0.1:34287
				Out: "unavailable: " + serverFail.URL,
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
			// netwait http://127.0.0.1:34287 http://127.0.0.1:34287
			cmd: icmd.Command(netwaitPath, server.URL, server.URL),
			expected: icmd.Expected{
				ExitCode: 0,
				// available: http://127.0.0.1:34287
				// available: http://127.0.0.1:34287
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
