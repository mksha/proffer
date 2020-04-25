package command

import (
	"example.com/proffer/components"
	awscopyamiresource "example.com/proffer/resources/aws/copyami"
	awsshareamiresource "example.com/proffer/resources/aws/shareami"
)

var Resources = map[string]components.Resourcer{
	"aws-copyami":  new(awscopyamiresource.Resource),
	"aws-shareami": new(awsshareamiresource.Resource),
}
