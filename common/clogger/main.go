package clogger

import (
	"fmt"
	"io"
	"log"
	"os"
)

var (
	LogLevel int // global log level.
)

type CLogger struct {
	Logger *log.Logger
}

// Log level.
const (
	DEBUG = iota
	INFO
	WARN
	ERROR
	PANIC
	FATAL
	NONE
)

// Base Style.
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

// Foreground text colors.
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

// Foreground Hi-Intensity text colors.
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

// Background text colors.
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

// Background Hi-Intensity text colors.
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
	escape    = "\x1b"
)


// New creates a custom logger with given args.
func New(out io.Writer, prefix string, flag int) *CLogger {
	return &CLogger{Logger: log.New(out, getColoredMsg(prefix, FgCyan), flag)}
}

// SetGlobalLogLevel sets the log level at global scope.
func SetGlobalLogLevel(level int) {
	LogLevel = level
}

// getColoredMsg returns a given string after applying given colorcode.
func getColoredMsg(msg string, codeList ...int) string {
	for _, code := range codeList {
		msg = fmt.Sprintf("%s[%dm%s\x1b[0m", escape, code, msg)
	}

	return msg
}

// SetPrefix sets the prefix for the given clogger.
func (cl *CLogger) SetPrefix(prefix string, code ...int) {
	if len(code) == 0 {
		cl.Logger.SetPrefix(getColoredMsg(prefix, FgCyan))
		return
	}

	cl.Logger.SetPrefix(getColoredMsg(prefix, code...))
}

// GetPrefix gets the current prefix of clogger.
func (cl *CLogger) GetPrefix() string {
	return cl.Logger.Prefix()
}

// Error is used to print info in red color.
// Message is only printed if the log level is set to lesser or equal to ERROR.
func (cl *CLogger) Error(a ...interface{}) {
	if LogLevel <= ERROR {
		msg := fmt.Sprint(a...)
		if err := cl.Logger.Output(callDepth, getColoredMsg(msg, FgRed)); err != nil {
			log.Fatalln(err)
		}
	}
}

// Errorf is used to print info in red color.
// Message is only printed if the log level is set to lesser or equal to ERROR.
func (cl *CLogger) Errorf(format string, a ...interface{}) {
	if LogLevel <= ERROR {
		msg := fmt.Sprintf(format, a...)
		if err := cl.Logger.Output(callDepth, getColoredMsg(msg, FgRed)); err != nil {
			log.Fatalln(err)
		}
	}
}

// Warn is used to print a formatted warning message in yellow color.
// Info is only printed if the log level is set to lesser or equal to WARN.
func (cl *CLogger) Warn(a ...interface{}) {
	if LogLevel <= WARN {
		msg := fmt.Sprint(a...)
		if err := cl.Logger.Output(callDepth, getColoredMsg(msg, Bold, FgYellow)); err != nil {
			log.Fatalln(err)
		}
	}
}

// Warn is used to print a formatted warning message in yellow color.
// Info is only printed if the log level is set to lesser or equal to WARN.
func (cl *CLogger) Warnf(format string, a ...interface{}) {
	if LogLevel <= WARN {
		msg := fmt.Sprintf(format, a...)
		if err := cl.Logger.Output(callDepth, getColoredMsg(msg, Bold, FgYellow)); err != nil {
			log.Fatalln(err)
		}
	}
}

// Success is used to print a success message in bold green color.
// Info is  printed regardless of log level.
func (cl *CLogger) Success(a ...interface{}) {
	msg := fmt.Sprint(a...)
	oldPrefix := cl.GetPrefix()
	cl.SetPrefix(oldPrefix, Bold)

	if err := cl.Logger.Output(callDepth, getColoredMsg(msg, Bold, FgHiGreen)); err != nil {
		log.Fatalln(err)
	}

	cl.SetPrefix(oldPrefix)
}

// Successf is used to print a formatted success message in bold green color.
// Info is printed regardless of log level.
func (cl *CLogger) Successf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	oldPrefix := cl.GetPrefix()
	cl.SetPrefix(oldPrefix, Bold)

	if err := cl.Logger.Output(callDepth, getColoredMsg(msg, Bold, FgHiGreen)); err != nil {
		log.Fatalln(err)
	}

	cl.SetPrefix(oldPrefix)
}

// Info is used to print message in green color.
// It can be used to print the message with normal severity.
// Message is only printed if the log level is set to lesser or equal to INFO.
func (cl *CLogger) Info(a ...interface{}) {
	if LogLevel <= INFO {
		msg := fmt.Sprint(a...)
		if err := cl.Logger.Output(callDepth, getColoredMsg(msg, FgGreen)); err != nil {
			log.Fatalln(err)
		}
	}
}

// Infof is used to print formatted info message in green color.
// It can be used to print the message with normal severity.
// Message is only printed if the log level is set to lesser or equal to INFO.
func (cl *CLogger) Infof(format string, a ...interface{}) {
	if LogLevel <= INFO {
		msg := fmt.Sprintf(format, a...)
		if err := cl.Logger.Output(callDepth, getColoredMsg(msg, FgGreen)); err != nil {
			log.Fatalln(err)
		}
	}
}

// Debug is used to print debug message in hi-intensity green color.
// Message is only printed if the log level is set to lesser or equal to DEBUG.
func (cl *CLogger) Debug(a ...interface{}) {
	if LogLevel <= DEBUG {
		msg := fmt.Sprint(a...)
		if err := cl.Logger.Output(callDepth, getColoredMsg(msg, FgHiGreen)); err != nil {
			log.Fatalln(err)
		}
	}
}

// Debugf is used to print formatted debug message in hi-intensity green color.
// Message is only printed if the log level is set to lesser or equal to DEBUG.
func (cl *CLogger) Debugf(format string, a ...interface{}) {
	if LogLevel <= DEBUG {
		msg := fmt.Sprintf(format, a...)
		if err := cl.Logger.Output(callDepth, getColoredMsg(msg, FgHiGreen)); err != nil {
			log.Fatalln(err)
		}
	}
}

// Panic is used to print panic message in bold red color.
// Message is only printed if the log level is set to lesser or equal to PANIC.
func (cl *CLogger) Panic(a ...interface{}) {
	if LogLevel <= PANIC {
		msg := fmt.Sprint(a...)
		if err := cl.Logger.Output(callDepth, getColoredMsg(msg, Bold, FgRed)); err != nil {
			log.Fatalln(err)
		}

		panic(msg)
	}
}

// Panicf is used to print formatted panic message in bold red color.
// Message is only printed if the log level is set to lesser or equal to PANIC.
func (cl *CLogger) Panicf(format string, a ...interface{}) {
	if LogLevel <= PANIC {
		msg := fmt.Sprintf(format, a...)
		if err := cl.Logger.Output(callDepth, getColoredMsg(msg, Bold, FgRed)); err != nil {
			log.Fatalln(err)
		}

		panic(msg)
	}
}

// Fatalf is used to print formatted fatal message in bold hi-intensity red color.
// Message is only printed if the log level is set to lesser or equal to FATAL.
// At the end it will call os.Exit(1).
func (cl *CLogger) Fatal(a ...interface{}) {
	if LogLevel <= FATAL {
		msg := fmt.Sprint(a...)
		if err := cl.Logger.Output(callDepth, getColoredMsg(msg, Bold, FgHiRed)); err != nil {
			log.Fatalln(err)
		}

		os.Exit(1)
	}
}

// Fatalf is used to print formatted fatal message in bold hi-intensity red color.
// Message is only printed if the log level is set to lesser or equal to FATAL.
// At the end it will call os.Exit(1).
func (cl *CLogger) Fatalf(format string, a ...interface{}) {
	if LogLevel <= FATAL {
		msg := fmt.Sprintf(format, a...)
		if err := cl.Logger.Output(callDepth, getColoredMsg(msg, Bold, FgHiRed)); err != nil {
			log.Fatalln(err)
		}

		os.Exit(1)
	}
}
