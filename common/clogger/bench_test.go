package clogger

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"testing"
)

func BenchmarkNew(b *testing.B) {
	// bench combined time to run through all test cases
	for n := range newTestCases {
		tc := newTestCases[n]

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// ignoring errors and results because we're just timing function execution
				out := &bytes.Buffer{}
				_ = New(out, tc.prefix, tc.flag)
			}
		})
	}
}

func BenchmarkSetGlobalLogLevel(b *testing.B) {
	// bench combined time to run through all test cases
	for n := range setGlobalLogLevelTestCases {
		tc := setGlobalLogLevelTestCases[n]

		for i := 0; i < b.N; i++ {
			SetGlobalLogLevel(tc.level)
		}
	}
}

func Benchmark_getColoredMsg(b *testing.B) {
	// bench combined time to run through all test cases
	for n := range getColoredMsgTestCases {
		tc := getColoredMsgTestCases[n]

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// ignoring errors and results because we're just timing function execution
				_ = getColoredMsg(tc.msg, tc.codeList...)
			}
		})
	}
}

func BenchmarkCLogger_SetPrefix(b *testing.B) {
	// bench combined time to run through all test cases
	for n := range cLoggerSetPrefixTestCases {
		tc := cLoggerSetPrefixTestCases[n]

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// ignoring errors and results because we're just timing function execution
				out := &bytes.Buffer{}
				cl := &CLogger{
					Logger: log.New(out, tc.prefix, log.Lmsgprefix),
				}
				cl.SetPrefix(tc.prefix, tc.code...)
			}
		})
	}
}

func BenchmarkCLogger_Error(b *testing.B) {
	// bench combined time to run through all test cases
	for n := range clogTestCases {
		tc := clogTestCases[n]

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// ignoring errors and results because we're just timing function execution
				out := &bytes.Buffer{}
				cl := &CLogger{
					Logger: log.New(out, "", log.Lmsgprefix),
				}
				cl.Error(tc.a...)
			}
		})
	}
}

func BenchmarkCLogger_Errorf(b *testing.B) {
	// bench combined time to run through all test cases
	for n := range clogTestCases {
		tc := clogTestCases[n]

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// ignoring errors and results because we're just timing function execution
				out := &bytes.Buffer{}
				cl := &CLogger{
					Logger: log.New(out, "", log.Lmsgprefix),
				}
				cl.Errorf("%v %v %v %v %v", tc.a...)
			}
		})
	}
}

func BenchmarkCLogger_Warn(b *testing.B) {
	// bench combined time to run through all test cases
	for n := range clogTestCases {
		tc := clogTestCases[n]

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// ignoring errors and results because we're just timing function execution
				out := &bytes.Buffer{}
				cl := &CLogger{
					Logger: log.New(out, "", log.Lmsgprefix),
				}
				cl.Warn(tc.a...)
			}
		})
	}
}

func BenchmarkCLogger_Warnf(b *testing.B) {
	// bench combined time to run through all test cases
	for n := range clogTestCases {
		tc := clogTestCases[n]

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// ignoring errors and results because we're just timing function execution
				out := &bytes.Buffer{}
				cl := &CLogger{
					Logger: log.New(out, "", log.Lmsgprefix),
				}
				cl.Warnf("%v %v %v %v %v", tc.a...)
			}
		})
	}
}

func BenchmarkCLogger_Success(b *testing.B) {
	// bench combined time to run through all test cases
	for n := range clogTestCases {
		tc := clogTestCases[n]

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// ignoring errors and results because we're just timing function execution
				out := &bytes.Buffer{}
				cl := &CLogger{
					Logger: log.New(out, "", log.Lmsgprefix),
				}
				cl.Success(tc.a...)
			}
		})
	}
}

func BenchmarkCLogger_Successf(b *testing.B) {
	// bench combined time to run through all test cases
	for n := range clogTestCases {
		tc := clogTestCases[n]

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// ignoring errors and results because we're just timing function execution
				out := &bytes.Buffer{}
				cl := &CLogger{
					Logger: log.New(out, "", log.Lmsgprefix),
				}
				cl.Successf("%v %v %v %v %v", tc.a...)
			}
		})
	}
}

func BenchmarkCLogger_Info(b *testing.B) {
	// bench combined time to run through all test cases
	for n := range clogTestCases {
		tc := clogTestCases[n]

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// ignoring errors and results because we're just timing function execution
				out := &bytes.Buffer{}
				cl := &CLogger{
					Logger: log.New(out, "", log.Lmsgprefix),
				}
				cl.Info(tc.a...)
			}
		})
	}
}

func BenchmarkCLogger_Infof(b *testing.B) {
	// bench combined time to run through all test cases
	for n := range clogTestCases {
		tc := clogTestCases[n]

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// ignoring errors and results because we're just timing function execution
				out := &bytes.Buffer{}
				cl := &CLogger{
					Logger: log.New(out, "", log.Lmsgprefix),
				}
				cl.Infof("%v %v %v %v %v", tc.a...)
			}
		})
	}
}

func BenchmarkCLogger_Debug(b *testing.B) {
	// bench combined time to run through all test cases
	for n := range clogTestCases {
		tc := clogTestCases[n]

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// ignoring errors and results because we're just timing function execution
				out := &bytes.Buffer{}
				cl := &CLogger{
					Logger: log.New(out, "", log.Lmsgprefix),
				}
				cl.Debug(tc.a...)
			}
		})
	}
}

func BenchmarkCLogger_Debugf(b *testing.B) {
	// bench combined time to run through all test cases
	for n := range clogTestCases {
		tc := clogTestCases[n]

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// ignoring errors and results because we're just timing function execution
				out := &bytes.Buffer{}
				cl := &CLogger{
					Logger: log.New(out, "", log.Lmsgprefix),
				}
				cl.Debugf("%v %v %v %v %v", tc.a...)
			}
		})
	}
}

func BenchmarkCLogger_Panic(b *testing.B) {
	// bench combined time to run through all test cases
	for n := range clogTestCases {
		tc := clogTestCases[n]

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// ignoring errors and results because we're just timing function execution
				out := &bytes.Buffer{}
				cl := &CLogger{
					Logger: log.New(out, "", log.Lmsgprefix),
				}
				cl.Panic(tc.a...)
			}
		})
	}
}

func BenchmarkCLogger_Panicf(b *testing.B) {
	// bench combined time to run through all test cases
	for n := range clogTestCases {
		tc := clogTestCases[n]

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// ignoring errors and results because we're just timing function execution
				out := &bytes.Buffer{}
				cl := &CLogger{
					Logger: log.New(out, "", log.Lmsgprefix),
				}
				cl.Panicf("%v %v %v %v %v", tc.a...)
			}
		})
	}
}

func BenchmarkCLogger_Fatal(b *testing.B) {
	// bench combined time to run through all test cases
	for n := range clogTestCases {
		tc := clogTestCases[n]

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// ignoring errors and results because we're just timing function execution
				if os.Getenv("TRIGGER_FATAL") == "1" {
					out := &bytes.Buffer{}
					cl := &CLogger{
						Logger: log.New(out, "", log.Lmsgprefix),
					}
					cl.Fatal(tc.a...)

					return
				}

				prog := os.Args[0]
				cmd := exec.Command(prog, "-test.run=TestCLogger_Fatal")

				cmd.Env = append(os.Environ(), "TRIGGER_FATAL=1")
				_ = cmd.Run()
			}
		})
	}
}

func BenchmarkCLogger_Fatalf(b *testing.B) {
	// bench combined time to run through all test cases
	for n := range clogTestCases {
		tc := clogTestCases[n]

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// ignoring errors and results because we're just timing function execution
				if os.Getenv("TRIGGER_FATAL") == "1" {
					out := &bytes.Buffer{}
					cl := &CLogger{
						Logger: log.New(out, "", log.Lmsgprefix),
					}
					cl.Fatalf("%v %v %v %v %v", tc.a...)

					return
				}

				prog := os.Args[0]
				cmd := exec.Command(prog, "-test.run=TestCLogger_Fatal")

				cmd.Env = append(os.Environ(), "TRIGGER_FATAL=1")
				_ = cmd.Run()
			}
		})
	}
}
