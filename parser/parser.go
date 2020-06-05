package parser

import (
	"log"
	"os"

	clog "github.com/proffer/common/clogger"
	"github.com/proffer/components"
)

type (
	// DynamicVars stores dynamic variables.
	DynamicVars map[string]interface{}
	// DefaultVars stores default variables.
	DefaultVars map[string]map[string]interface{}
)

var (
	clogger     = clog.New(os.Stdout, "config-parser | ", log.Lmsgprefix)
	dynamicVars DynamicVars
	defaultVars DefaultVars
)

// TemplateConfig represent the proffer template file in go data-structure.
type TemplateConfig struct {
	Vars         map[string]interface{}   `yaml:"vars"`
	RawResources []components.RawResource `yaml:"resources,flow" required:"true"`
	Resources    components.MapOfResource `yaml:"-"`
}
