package clogger

import (
	"log"
)

// test cases for New function.
var newTestCases = []struct {
	name   string
	prefix string
	flag   int
}{
	{
		name:   "prefix with single flag",
		prefix: "Test:",
		flag:   log.Lmsgprefix,
	},
	{
		name:   "prefix with multiple flag",
		prefix: "Check:",
		flag:   log.Lmsgprefix | log.LstdFlags,
	},
	{
		name:   "empty prefix  with multiple flag",
		prefix: "",
		flag:   log.Lmicroseconds | log.LstdFlags,
	},
}

// test cases for SetGlobalLogLevel function.
var setGlobalLogLevelTestCases = []struct {
	name  string
	level int
}{
	{
		name:  "default log level",
		level: DEBUG,
	},
	{
		name:  "Set global log level Info",
		level: INFO,
	},
	{
		name:  "Set global log level warn",
		level: WARN,
	},
	{
		name:  "Set global log level error",
		level: ERROR,
	},
	{
		name:  "Set global log level panic",
		level: PANIC,
	},
	{
		name:  "Set global log level fatal",
		level: FATAL,
	},
	{
		name:  "disable logging",
		level: NONE,
	},
}

// test cases for getColoredMsg function.
var getColoredMsgTestCases = []struct {
	name     string
	msg      string
	codeList []int
	want     string
}{
	{
		name:     "msg with color code",
		msg:      "Hello",
		codeList: []int{FgRed},
		want:     "\x1b[31mHello\x1b[0m",
	},
	{
		name:     "msg with multiple color code",
		msg:      "Hello",
		codeList: []int{FgRed, Bold},
		want:     "\x1b[1m\x1b[31mHello\x1b[0m\x1b[0m",
	},
}

// test cases for CLogger_SetPrefix function.
var cLoggerSetPrefixTestCases = []struct {
	name   string
	prefix string
	code   []int
}{
	{
		name:   "prefix with single code",
		prefix: "Hello",
		code:   []int{FgBlue},
	},
	{
		name:   "prefix with multiple code",
		prefix: "Hello",
		code:   []int{FgBlue, Bold},
	},
}

// test cases for diff clog functions.
var clogTestCases = []struct {
	name string
	a    []interface{}
}{
	{
		name: "call with no argument",
		a:    []interface{}{},
	},
	{
		name: "call with multiple arguments of diff types",
		a:    []interface{}{"Test", 1, 4.3, map[string]string{"Hello": "Hi"}, struct{ Name string }{Name: "Mohit"}},
	},
}
