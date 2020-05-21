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
	Rmicroseconds = `\.[0-9][0-9][0-9][0-9][0-9][0-9]`
)

func TestNew(t *testing.T) {
	tests := []struct {
		name   string
		prefix string
		flag   int
		want   string
	}{
		{
			name:   "prefix with single flag",
			prefix: "Test:",
			flag:   log.Lmsgprefix,
			want:   "Test:*",
		},
		{
			name:   "prefix with multiple flag",
			prefix: "Check:",
			flag:   log.Lmsgprefix | log.LstdFlags,
			want:   Rdate + " " + Rtime + " " + "Check:*",
		},
		{
			name:   "empty prefix  with multiple flag",
			prefix: "",
			flag:   log.Lmicroseconds | log.LstdFlags,
			want:   Rdate + " " + Rmicroseconds + " *",
		},
	}

	for n := range tests {
		tt := tests[n]
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			cl := New(out, tt.prefix, tt.flag)
			gotP := cl.GetPrefix()
			wantP := fmt.Sprintf("%s[%dm%s\x1b[0m", escape, FgCyan, tt.prefix)

			if gotP != wantP {
				t.Errorf("Prefix: expected %s, got %s", wantP, gotP)
			}

			// cl.Logger.Print("Hello")
			// if
			// fmt.Println(out.String())
			// fmt.Println(regexp.Match(tt.want, out.Bytes()))
		})
	}
}
