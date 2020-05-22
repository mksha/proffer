package clogger

import (
	"bytes"
	"fmt"
	"log"
	"testing"
)

const (
	Rdate         = `[0-9][0-9][0-9][0-9]/[0-9][0-9]/[0-9][0-9]`
	Rtime         = `[0-9][0-9]:[0-9][0-9]:[0-9][0-9]`
	Rmicroseconds = Rtime + `.[0-9][0-9][0-9][0-9][0-9][0-9]`
)

func TestNew(t *testing.T) {
	for n := range newTestCases {
		tt := newTestCases[n]
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			cl := New(out, tt.prefix, tt.flag)
			gotP := cl.GetPrefix()
			wantP := fmt.Sprintf("%s[%dm%s\x1b[0m", escape, FgCyan, tt.prefix)

			if gotP != wantP {
				t.Errorf("Prefix: expected %s, got %s", wantP, gotP)
			}
		})
	}
}

func TestSetGlobalLogLevel(t *testing.T) {
	for n := range setGlobalLogLevelTestCases {
		tt := setGlobalLogLevelTestCases[n]
		t.Run(tt.name, func(t *testing.T) {
			SetGlobalLogLevel(tt.level)
			if LogLevel != tt.level {
				t.Errorf("want: %v , got: %v", tt.level, LogLevel)
			}
		})
	}
}

func Test_getColoredMsg(t *testing.T) {
	for n := range getColoredMsgTestCases {
		tt := getColoredMsgTestCases[n]
		t.Run(tt.name, func(t *testing.T) {
			if got := getColoredMsg(tt.msg, tt.codeList...); got != tt.want {
				t.Errorf("getColoredMsg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCLogger_SetPrefix(t *testing.T) {
	for n := range cLoggerSetPrefixTestCases {
		tt := cLoggerSetPrefixTestCases[n]

		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			cl := &CLogger{
				Logger: log.New(out, tt.prefix, log.Lmsgprefix),
			}
			cl.SetPrefix(tt.prefix, tt.code...)
			want := getColoredMsg(tt.prefix, tt.code...)
			got := cl.GetPrefix()
			if got != want {
				t.Errorf("want: %v, got: %v", want, got)
			}
		})
	}
}

func TestCLogger_Error(t *testing.T) {
	for n := range clogTestCases {
		tt := clogTestCases[n]
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			cl := &CLogger{
				Logger: log.New(out, "", log.Lmsgprefix),
			}

			if LogLevel > ERROR {
				cl.Error(tt.a...)

				want := ""
				got := out.String()
				if want != got {
					t.Errorf("want: %v, got: %v", want, got)
				}
			}

			LogLevel = ERROR

			cl.Error(tt.a...)

			got := out.String()
			msg := fmt.Sprint(tt.a...)
			want := fmt.Sprintf("%s[%dm%s\x1b[0m\n", escape, FgRed, msg)

			if want != got {
				t.Errorf("want: %v, got: %v", want, got)
			}
		})
	}
}

func TestCLogger_Errorf(t *testing.T) {
	for n := range clogTestCases {
		tt := clogTestCases[n]
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			cl := &CLogger{
				Logger: log.New(out, "", log.Lmsgprefix),
			}

			if LogLevel > ERROR {
				cl.Errorf("%s %d %f %v %v", tt.a...)

				got := out.String()
				want := ""

				if want != got {
					t.Errorf("want: %v, got: %v", want, got)
				}
			}

			LogLevel = ERROR

			cl.Errorf("%s %d %f %v %v", tt.a...)

			got := out.String()
			msg := fmt.Sprintf("%s %d %f %v %v", tt.a...)
			want := fmt.Sprintf("%s[%dm%s\x1b[0m\n", escape, FgRed, msg)

			if want != got {
				t.Errorf("want: %v, got: %v", want, got)
			}
		})
	}
}

func TestCLogger_Warn(t *testing.T) {
	for n := range clogTestCases {
		tt := clogTestCases[n]
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			cl := &CLogger{
				Logger: log.New(out, "", log.Lmsgprefix),
			}

			if LogLevel > WARN {
				cl.Warn(tt.a...)

				want := ""
				got := out.String()

				if want != got {
					t.Errorf("want: %v, got: %v", want, got)
				}
			}

			LogLevel = WARN

			cl.Warn(tt.a...)

			got := out.String()
			msg := fmt.Sprint(tt.a...)
			msg = fmt.Sprintf("%s[%dm%s\x1b[0m", escape, Bold, msg)
			want := fmt.Sprintf("%s[%dm%s\x1b[0m\n", escape, FgYellow, msg)

			if want != got {
				t.Errorf("want: %v, got: %v", want, got)
			}
		})
	}
}

func TestCLogger_Warnf(t *testing.T) {
	for n := range clogTestCases {
		tt := clogTestCases[n]
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			cl := &CLogger{
				Logger: log.New(out, "", log.Lmsgprefix),
			}

			if LogLevel > WARN {
				cl.Warnf("%s %d %f %v %v", tt.a...)

				got := out.String()
				want := ""

				if want != got {
					t.Errorf("want: %v, got: %v", want, got)
				}
			}

			LogLevel = WARN

			cl.Warnf("%s %d %f %v %v", tt.a...)

			got := out.String()
			msg := fmt.Sprintf("%s %d %f %v %v", tt.a...)
			msg = fmt.Sprintf("%s[%dm%s\x1b[0m", escape, Bold, msg)
			want := fmt.Sprintf("%s[%dm%s\x1b[0m\n", escape, FgYellow, msg)

			if want != got {
				t.Errorf("want: %v, got: %v", want, got)
			}
		})
	}
}

// func TestCLogger_Success(t *testing.T) {
// 	for n := range clogTestCases {
// 		tt := clogTestCases[n]
// 		t.Run(tt.name, func(t *testing.T) {
// 			out := &bytes.Buffer{}
// 			cl := &CLogger{
// 				Logger: log.New(out, "", log.Lmsgprefix),
// 			}

// 			cl.Success(tt.a...)

// 			got := out.String()
// 			msg := fmt.Sprint(tt.a...)
// 			pre := fmt.Sprintf("%s[%dm%s\x1b[0m", escape, FgCyan, "")
// 			pre = fmt.Sprintf("%s[%dm%s\x1b[0m", escape, Bold, pre)
// 			msg = fmt.Sprintf("%s[%dm%s\x1b[0m", escape, Bold, pre+msg)
// 			want := fmt.Sprintf("%s[%dm%s\x1b[0m", escape, FgHiGreen, msg)

// 			fmt.Println(len(want))
// 			fmt.Println(len(got))
// 			if want != got {
// 				t.Errorf("want: %v, got: %v", want, got)
// 			}
// 		})
// 	}
// }

func TestCLogger_Info(t *testing.T) {
	for n := range clogTestCases {
		tt := clogTestCases[n]
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			cl := &CLogger{
				Logger: log.New(out, "", log.Lmsgprefix),
			}

			if LogLevel > INFO {
				cl.Info(tt.a...)

				want := ""
				got := out.String()

				if want != got {
					t.Errorf("want: %v, got: %v", want, got)
				}
			}

			LogLevel = INFO

			cl.Info(tt.a...)

			got := out.String()
			msg := fmt.Sprint(tt.a...)
			want := fmt.Sprintf("%s[%dm%s\x1b[0m\n", escape, FgGreen, msg)

			if want != got {
				t.Errorf("want: %v, got: %v", want, got)
			}
		})
	}
}

func TestCLogger_Infof(t *testing.T) {
	for n := range clogTestCases {
		tt := clogTestCases[n]
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			cl := &CLogger{
				Logger: log.New(out, "", log.Lmsgprefix),
			}

			if LogLevel > INFO {
				cl.Infof("%s %d %f %v %v", tt.a...)

				got := out.String()
				want := ""

				if want != got {
					t.Errorf("want: %v, got: %v", want, got)
				}
			}

			LogLevel = INFO

			cl.Infof("%s %d %f %v %v", tt.a...)

			got := out.String()
			msg := fmt.Sprintf("%s %d %f %v %v", tt.a...)
			want := fmt.Sprintf("%s[%dm%s\x1b[0m\n", escape, FgGreen, msg)

			if want != got {
				t.Errorf("want: %v, got: %v", want, got)
			}
		})
	}
}

func TestCLogger_Debug(t *testing.T) {
	for n := range clogTestCases {
		tt := clogTestCases[n]
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			cl := &CLogger{
				Logger: log.New(out, "", log.Lmsgprefix),
			}

			if LogLevel > DEBUG {
				cl.Debug(tt.a...)

				want := ""
				got := out.String()

				if want != got {
					t.Errorf("want: %v, got: %v", want, got)
				}
			}

			LogLevel = DEBUG

			cl.Debug(tt.a...)

			got := out.String()
			msg := fmt.Sprint(tt.a...)
			want := fmt.Sprintf("%s[%dm%s\x1b[0m\n", escape, FgHiGreen, msg)

			if want != got {
				t.Errorf("want: %v, got: %v", want, got)
			}
		})
	}
}

func TestCLogger_Debugf(t *testing.T) {
	for n := range clogTestCases {
		tt := clogTestCases[n]
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			cl := &CLogger{
				Logger: log.New(out, "", log.Lmsgprefix),
			}

			if LogLevel > DEBUG {
				cl.Debugf("%s %d %f %v %v", tt.a...)

				got := out.String()
				want := ""

				if want != got {
					t.Errorf("want: %v, got: %v", want, got)
				}
			}

			LogLevel = DEBUG

			cl.Debugf("%s %d %f %v %v", tt.a...)

			got := out.String()
			msg := fmt.Sprintf("%s %d %f %v %v", tt.a...)
			want := fmt.Sprintf("%s[%dm%s\x1b[0m\n", escape, FgHiGreen, msg)

			if want != got {
				t.Errorf("want: %v, got: %v", want, got)
			}
		})
	}
}

func TestCLogger_Panic(t *testing.T) {
	for n := range clogTestCases {
		tt := clogTestCases[n]
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Error("Panic not triggered")
				}
			}()
			out := &bytes.Buffer{}
			cl := &CLogger{
				Logger: log.New(out, "", log.Lmsgprefix),
			}
			cl.Panic(tt.a...)
		})
	}
}

func TestCLogger_Panicf(t *testing.T) {
	for n := range clogTestCases {
		tt := clogTestCases[n]
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Error("Panic not triggered")
				}
			}()
			out := &bytes.Buffer{}
			cl := &CLogger{
				Logger: log.New(out, "", log.Lmsgprefix),
			}
			cl.Panicf("%s %d %v %v %v", tt.a...)
		})
	}
}

// func TestCLogger_Fatal(t *testing.T) {
// 	for n := range clogTestCases {
// 		tt := clogTestCases[n]
// 		t.Run(tt.name, func(t *testing.T) {
// 			out := &bytes.Buffer{}
// 			cl := &CLogger{
// 				Logger: log.New(out, "", log.Lmsgprefix),
// 			}
// 			var flag bool
// 			go func (){
// 				cl.Fatal(tt.a...)
// 				flag = true
// 			}()

// 			if flag {
// 				t.Error("Did not cause exit")
// 			}
// 		})
// 	}
// }
