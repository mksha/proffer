package parser

import (
	"log"
	"os"

	clog "example.com/proffer/common/clogger"
)

var (
	clogger = clog.New(os.Stdout, "config-parser | ", log.Lmsgprefix)
)
