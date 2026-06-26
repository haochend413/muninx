//go:build !windows

package cmd

import (
	"os/exec"
	"syscall"
)

// detachProcess puts the child in its own process group so it doesn't receive
// signals (e.g. SIGINT from Ctrl+C) that are sent to muninx's process group.
func detachProcess(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}
