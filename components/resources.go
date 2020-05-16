package components

type Resourcer interface {
	Validate(RawResource) error
	Prepare(map[string]interface{}) error
	Run() error
}

type RawResource struct {
	Name   string                 `yaml:"name" required:"true"`
	Type   string                 `yaml:"type" required:"true"`
	Config map[string]interface{} `yaml:"config" required:"true"`
}

type MapOfResource map[string]func() (Resourcer, error)
