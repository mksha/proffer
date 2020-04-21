package command

import (
	"example.com/amidist/components"
	awscopyamiresource "example.com/amidist/resources/aws/copyami"
	awsshareamiresource "example.com/amidist/resources/aws/shareami"
)

var Resources = map[string]components.Resourcer{
	"aws-copyami":  new(awscopyamiresource.Resource),
	"aws-shareami": new(awsshareamiresource.Resource),
}
