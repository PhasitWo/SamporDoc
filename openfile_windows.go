//go:build windows

package main

import (
	"os/exec"
	"syscall"
)

func (a *App) CMDOpenFile(filePath string) error {
	var cmd *exec.Cmd
	cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", filePath)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000,
	}
	return cmd.Run()
}
