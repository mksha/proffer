// +build windows

package clogger

import (
	"os"
	"runtime"

	"golang.org/x/sys/windows"
)

func init() {
	// Enable color logging for windows platform
	if runtime.GOOS == "windows" {
		stdout := windows.Handle(os.Stdout.Fd())
		var originalMode uint32

		windows.GetConsoleMode(stdout, &originalMode)
		windows.SetConsoleMode(stdout, originalMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
	}
}
