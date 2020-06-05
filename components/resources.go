package components

// Resourcer implements identity of resource type.
type Resourcer interface {
	Validate(RawResource) error
	Prepare(map[string]interface{}) error
	Run() error
}

// RawResource represents raw resource information.
type RawResource struct {
	Name   string                 `yaml:"name" required:"true"`
	Type   string                 `yaml:"type" required:"true"`
	Config map[string]interface{} `yaml:"config" required:"true"`
}

// MapOfResource represents map of resources.
type MapOfResource map[string]func() (Resourcer, error)
