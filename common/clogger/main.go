package clogger

import (
	"fmt"
	"io"
	"log"
	"os"
)

var (
	LogLevel int
)

type CLogger struct {
	Logger *log.Logger
}

// Log level
const (
	DEBUG = iota
	INFO
	WARN
	ERROR
	PANIC
	FATAL
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

func New(out io.Writer, prefix string, flag int) *CLogger {
	return &CLogger{Logger: log.New(out, getColoredMsg(prefix, FgCyan), flag)}
}
func SetGlobalLogLevel(level int) {
	LogLevel = level
}

func getColoredMsg(msg string, codeList ...int) string {
	for _, code := range codeList {
		msg = fmt.Sprintf("\x1b[%dm%s\x1b[0m", code, msg)
	}

	return msg
}

func getFormattedMsg(format string, a ...interface{}) string {
	if len(a) == 0 {
		return format
	}

	return fmt.Sprintf(format, a...)
}

func getMsg(a ...interface{}) string {
	return fmt.Sprint(a...)
}

func (cl *CLogger) SetPrefix(prefix string) {
	cl.Logger.SetPrefix(getColoredMsg(prefix, FgCyan))
}

// Error is used to print info in red color
func (cl *CLogger) Error(a ...interface{}) {
	if LogLevel <= ERROR {
		if err := cl.Logger.Output(2, getColoredMsg(getMsg(a...), Bold, FgRed)); err != nil {
			log.Fatalln(err)
		}
	}
}

// Errorf is used to print info in red color
func (cl *CLogger) Errorf(format string, a ...interface{}) {
	if LogLevel <= ERROR {
		if err := cl.Logger.Output(2, getColoredMsg(getFormattedMsg(format, a...), Bold, FgRed)); err != nil {
			log.Fatalln(err)
		}
	}
}

// Warn is used to print a formatted warning message in yellow color
func (cl *CLogger) Warn(a ...interface{}) {
	if LogLevel <= WARN {
		if err := cl.Logger.Output(2, getColoredMsg(getMsg(a...), Bold, FgYellow)); err != nil {
			log.Fatalln(err)
		}
	}
}

// Warn is used to print a warning message in yellow color
func (cl *CLogger) Warnf(format string, a ...interface{}) {
	if LogLevel <= WARN {
		if err := cl.Logger.Output(2, getColoredMsg(getFormattedMsg(format, a...), Bold, FgYellow)); err != nil {
			log.Fatalln(err)
		}
	}
}

// Success is used to print a success message in green color
func Success(msg string) {
	fmt.Println(getColoredMsg(msg, Bold, FgGreen))
}

// Successf is used to print a success message in green color
func Successf(format string, a ...interface{}) {
	fmt.Println(getColoredMsg(getFormattedMsg(format, a...), Bold, FgGreen))
}

// Info is used to print info message in blue color

func (cl *CLogger) Info(a ...interface{}) {
	if LogLevel <= INFO {
		if err := cl.Logger.Output(2, getColoredMsg(getMsg(a...), FgGreen)); err != nil {
			log.Fatalln(err)
		}
	}
}

// Infof is used to print formatted info message in blue color
func (cl *CLogger) Infof(format string, a ...interface{}) {
	if LogLevel <= INFO {
		if err := cl.Logger.Output(2, getColoredMsg(getFormattedMsg(format, a...), FgGreen)); err != nil {
			log.Fatalln(err)
		}
	}
}

// Debug is used to print debug message in green color
func (cl *CLogger) Debug(a ...interface{}) {
	if LogLevel <= DEBUG {
		if err := cl.Logger.Output(2, getColoredMsg(getMsg(a...), Bold, FgGreen)); err != nil {
			log.Fatalln(err)
		}
	}
}

// Debugf is used to print formatted debug message in green color
func (cl *CLogger) Debugf(format string, a ...interface{}) {
	if LogLevel <= DEBUG {
		if err := cl.Logger.Output(2, getColoredMsg(getFormattedMsg(format, a...), Bold, FgGreen)); err != nil {
			log.Fatalln(err)
		}
	}
}

// Fatal is used to print debug message in green color
func (cl *CLogger) Fatal(a ...interface{}) {
	if LogLevel <= FATAL {
		if err := cl.Logger.Output(2, getColoredMsg(getMsg(a...), Bold, FgHiRed)); err != nil {
			log.Fatalln(err)
		}
		os.Exit(1)
	}
}

// Fatalf is used to print formatted debug message in green color
func (cl *CLogger) Fatalf(format string, a ...interface{}) {
	if LogLevel <= DEBUG {
		if err := cl.Logger.Output(2, getColoredMsg(getFormattedMsg(format, a...), Bold, FgHiRed)); err != nil {
			log.Fatalln(err)
		}
		os.Exit(1)
	}
}
