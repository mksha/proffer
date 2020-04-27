package common

import (
	"fmt"
)

// Base Style
const (
	Reset = iota
	Bold
	Faint
	Italic
	Underline
	BlinkSlow
	BlinkRapid
	ReverseVideo
	Concealed
	CrossedOut
)

// Foreground text colors
const (
	FgBlack = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
)

// Foreground Hi-Intensity text colors
const (
	FgHiBlack = iota + 90
	FgHiRed
	FgHiGreen
	FgHiYellow
	FgHiBlue
	FgHiMagenta
	FgHiCyan
	FgHiWhite
)

// Background text colors
const (
	BgBlack = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
)

// Background Hi-Intensity text colors
const (
	BgHiBlack = iota + 100
	BgHiRed
	BgHiGreen
	BgHiYellow
	BgHiBlue
	BgHiMagenta
	BgHiCyan
	BgHiWhite
)

func getColoredMsg(msg string, codeList ...int) string {
	for _, code := range codeList {
		msg = fmt.Sprintf("\x1b[%dm%s\x1b[0m", code, msg)
	}

	return msg
}

func getFormattedMsg(format string, a ...interface{}) string {
	if len(a) == 0 {
		return fmt.Sprintf(format)
	}

	return fmt.Sprintf(format, a...)
}

// Error is used to print info in red color
func Error(format string, a ...interface{}) string {
	return getColoredMsg(getFormattedMsg(format, a...), Bold, FgRed)
}

// Warning is used to print a warning message in yellow color
func Warning(format string, a ...interface{}) string {
	return getColoredMsg(getFormattedMsg(format, a...), Bold, FgYellow)
}

// Success is used to print a success message in green color

func Success(format string, a ...interface{}) string {
	return getColoredMsg(getFormattedMsg(format, a...), Bold, FgGreen)
}

// Info is used to print info message in blue color
func Info(format string, a ...interface{}) string {
	return getColoredMsg(getFormattedMsg(format, a...), Bold, FgBlue)
}
