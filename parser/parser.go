package parser

import (
	"log"
	"os"

	clog "example.com/proffer/common/clogger"
	"example.com/proffer/components"
)

var (
	clogger = clog.New(os.Stdout, "config-parser | ", log.Lmsgprefix)
)

type TemplateConfig struct {
	Variables    map[string]string        `yaml:"variables,flow"`
	RawResources []components.RawResource `yaml:"resources,flow" required:"true"`
	Resources    components.MapOfResource `yaml:"-"`
}
