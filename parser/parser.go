package parser

import (
	"log"
	"os"

	clog "github.com/proffer/common/clogger"
	"github.com/proffer/components"
)

type (
	DynamicVars map[string]interface{}
	DefaultVars map[string]map[string]interface{}
)

var (
	clogger     = clog.New(os.Stdout, "config-parser | ", log.Lmsgprefix)
	dynamicVars DynamicVars
	defaultVars DefaultVars
)

type TemplateConfig struct {
	Vars         map[string]interface{}   `yaml:"vars"`
	RawResources []components.RawResource `yaml:"resources,flow" required:"true"`
	Resources    components.MapOfResource `yaml:"-"`
}
