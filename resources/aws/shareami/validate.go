package shareami

import (
	"fmt"

	"example.com/proffer/components"
	"github.com/mitchellh/mapstructure"
)

func (r *Resource) Validate(rawResource components.RawResource) error {
	var c Config

	r.Name = &rawResource.Name
	r.Type = &rawResource.Type

	var md mapstructure.Metadata

	if err := mapstructure.DecodeMetadata(rawResource.Config, &c, &md); err != nil {
		return err
	}

	if len(md.Unused) != 0 {
		clogger.Errorf("Invalid key(s) in config section of Resource: [%s]", *r.Name)
		clogger.Error(md.Unused)

		return fmt.Errorf("validation failed")
	}

	return nil
}
