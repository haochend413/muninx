package cmd

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

const embedBaseURL = "http://127.0.0.1:8000"

// embedServerRunning returns true if the embedding server responds on /docs.
func embedServerRunning() bool {
	c := &http.Client{Timeout: 2 * time.Second}
	resp, err := c.Get(embedBaseURL + "/docs")
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode == 200
}

// launchEmbedServer starts the uvicorn embedding server in the background,
// mirroring what scripts/build.sh does. It blocks until the server is ready
// (up to 5 seconds) and returns an error if it never comes up.
func launchEmbedServer(venvDir, serverDir string) error {
	if embedServerRunning() {
		fmt.Println("Embedding server already running.")
		return nil
	}

	if venvDir == "" {
		return fmt.Errorf("VenvDir is not set in config; cannot start embedding server")
	}
	if serverDir == "" {
		return fmt.Errorf("EmbedServerDir is not set in config; cannot start embedding server")
	}

	var pythonBin string
	if runtime.GOOS == "windows" {
		pythonBin = filepath.Join(venvDir, "Scripts", "python.exe")
	} else {
		pythonBin = filepath.Join(venvDir, "bin", "python")
	}

	if _, err := os.Stat(pythonBin); err != nil {
		return fmt.Errorf("python not found at %s — check VenvDir in config", pythonBin)
	}
	if _, err := os.Stat(serverDir); err != nil {
		return fmt.Errorf("server dir not found at %s — check EmbedServerDir in config", serverDir)
	}

	logPath := filepath.Join(serverDir, "embed_server.log")
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("cannot open embed server log: %w", err)
	}
	defer logFile.Close()

	proc := exec.Command(pythonBin, "-m", "uvicorn", "embed:app",
		"--app-dir", serverDir,
		"--host", "127.0.0.1",
		"--port", "8000",
	)
	proc.Stdout = logFile
	proc.Stderr = logFile
	detachProcess(proc)

	if err := proc.Start(); err != nil {
		return fmt.Errorf("failed to start embedding server: %w", err)
	}
	fmt.Printf("Starting embedding server (PID %d), logs: %s\n", proc.Process.Pid, logPath)

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		time.Sleep(500 * time.Millisecond)
		if embedServerRunning() {
			fmt.Println("Embedding server ready.")
			return nil
		}
	}
	return fmt.Errorf("embedding server did not become ready in time; check logs: %s", logPath)
}
