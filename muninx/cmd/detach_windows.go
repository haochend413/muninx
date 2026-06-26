//go:build windows

package cmd

import "os/exec"

func detachProcess(cmd *exec.Cmd) {}
