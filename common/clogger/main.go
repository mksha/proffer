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

const (
	callDepth = 2
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

func (cl *CLogger) SetPrefix(prefix string) {
	cl.Logger.SetPrefix(getColoredMsg(prefix, FgCyan))
}

// Error is used to print info in red color
func (cl *CLogger) Error(a ...interface{}) {
	if LogLevel <= ERROR {
		msg := fmt.Sprint(a...)
		if err := cl.Logger.Output(callDepth, getColoredMsg(msg, FgRed)); err != nil {
			log.Fatalln(err)
		}
	}
}

// Errorf is used to print info in red color
func (cl *CLogger) Errorf(format string, a ...interface{}) {
	if LogLevel <= ERROR {
		msg := fmt.Sprintf(format, a...)
		if err := cl.Logger.Output(callDepth, getColoredMsg(msg, FgRed)); err != nil {
			log.Fatalln(err)
		}
	}
}

// Warn is used to print a formatted warning message in yellow color
func (cl *CLogger) Warn(a ...interface{}) {
	if LogLevel <= WARN {
		msg := fmt.Sprint(a...)
		if err := cl.Logger.Output(callDepth, getColoredMsg(msg, Bold, FgYellow)); err != nil {
			log.Fatalln(err)
		}
	}
}

// Warn is used to print a warning message in yellow color
func (cl *CLogger) Warnf(format string, a ...interface{}) {
	if LogLevel <= WARN {
		msg := fmt.Sprintf(format, a...)
		if err := cl.Logger.Output(callDepth, getColoredMsg(msg, Bold, FgYellow)); err != nil {
			log.Fatalln(err)
		}
	}
}

// Success is used to print a success message in green color
func (cl *CLogger) Success(a ...interface{}) {
	msg := fmt.Sprint(a...)
	if err := cl.Logger.Output(callDepth, getColoredMsg(msg, Bold, FgHiGreen)); err != nil {
		log.Fatalln(err)
	}
}

// Successf is used to print a success message in green color
func (cl *CLogger) Successf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	if err := cl.Logger.Output(callDepth, getColoredMsg(msg, Bold, FgHiGreen)); err != nil {
		log.Fatalln(err)
	}
}

// Info is used to print info message in blue color

func (cl *CLogger) Info(a ...interface{}) {
	if LogLevel <= INFO {
		msg := fmt.Sprint(a...)
		if err := cl.Logger.Output(callDepth, getColoredMsg(msg, FgGreen)); err != nil {
			log.Fatalln(err)
		}
	}
}

// Infof is used to print formatted info message in blue color
func (cl *CLogger) Infof(format string, a ...interface{}) {
	if LogLevel <= INFO {
		msg := fmt.Sprintf(format, a...)
		if err := cl.Logger.Output(callDepth, getColoredMsg(msg, FgGreen)); err != nil {
			log.Fatalln(err)
		}
	}
}

// Debug is used to print debug message in green color
func (cl *CLogger) Debug(a ...interface{}) {
	if LogLevel <= DEBUG {
		msg := fmt.Sprint(a...)
		if err := cl.Logger.Output(callDepth, getColoredMsg(msg, FgHiGreen)); err != nil {
			log.Fatalln(err)
		}
	}
}

// Debugf is used to print formatted debug message in green color
func (cl *CLogger) Debugf(format string, a ...interface{}) {
	if LogLevel <= DEBUG {
		msg := fmt.Sprintf(format, a...)
		if err := cl.Logger.Output(callDepth, getColoredMsg(msg, FgHiGreen)); err != nil {
			log.Fatalln(err)
		}
	}
}

// Panic is used to print panic message in green color
func (cl *CLogger) Panic(a ...interface{}) {
	if LogLevel <= PANIC {
		msg := fmt.Sprint(a...)
		if err := cl.Logger.Output(callDepth, getColoredMsg(msg, Bold, FgRed)); err != nil {
			log.Fatalln(err)
		}

		panic(msg)
	}
}

// Panicf is used to print formatted panic message in green color
func (cl *CLogger) Panicf(format string, a ...interface{}) {
	if LogLevel <= PANIC {
		msg := fmt.Sprintf(format, a...)
		if err := cl.Logger.Output(callDepth, getColoredMsg(msg, Bold, FgRed)); err != nil {
			log.Fatalln(err)
		}

		panic(msg)
	}
}

// Fatalf is used to print formatted fatal message in green color
func (cl *CLogger) Fatal(a ...interface{}) {
	if LogLevel <= FATAL {
		msg := fmt.Sprint(a...)
		if err := cl.Logger.Output(callDepth, getColoredMsg(msg, Bold, FgHiRed)); err != nil {
			log.Fatalln(err)
		}

		os.Exit(1)
	}
}

// Fatalf is used to print formatted fatal message in green color
func (cl *CLogger) Fatalf(format string, a ...interface{}) {
	if LogLevel <= FATAL {
		msg := fmt.Sprintf(format, a...)
		if err := cl.Logger.Output(callDepth, getColoredMsg(msg, Bold, FgHiRed)); err != nil {
			log.Fatalln(err)
		}

		os.Exit(1)
	}
}
