//go:build !windows

package main

import (
	"os/exec"
	"runtime"
)

func (a *App) CMDOpenFile(filePath string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", filePath)
	default:
		// Fallback for other systems, e.g., Linux
		cmd = exec.Command("xdg-open", filePath)
	}
	return cmd.Run()
}
