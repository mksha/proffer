package shareami

import (
	"github.com/aws/aws-sdk-go/service/ec2"
)

type TargetInfo struct {
	Regions  []*string
	CopyTags bool
	Tags     []*ec2.Tag
}

func shareAmi() error {
	return nil
}
