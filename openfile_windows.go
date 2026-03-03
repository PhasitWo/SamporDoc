//go:build windows

package main

import (
	"golang.org/x/sys/windows"
)

func (a *App) CMDOpenFile(filePath string) error {
	verbPtr, _ := windows.UTF16PtrFromString("open")
	pathPtr, _ := windows.UTF16PtrFromString(filePath)

	err := windows.ShellExecute(
		0,                     // Parent window handle (0 for none)
		verbPtr,               // The action to perform
		pathPtr,               // The file path
		nil,                   // Parameters (used if opening an .exe)
		nil,                   // Working directory (nil uses current)
		windows.SW_SHOWNORMAL, // How to show the window
	)

	return err
}
