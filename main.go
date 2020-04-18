package main

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {
	cmdArgs := os.Args
	dsc := cmdArgs[1]
	parseTmpl(dsc)
	srcAmiInfo := SrcAmiInfo{
		Region: aws.String("us-west-2"),
		Filters: []*ec2.Filter{
			{
				Name: aws.String("image-id"),
				Values: []*string{
					aws.String("ami-01b55974057720a43"),
				},
			},
		},
	}

	targetInfo := TargetInfo{
		Regions: []*string{
			aws.String("us-west-1"),
		},
	}
	copyAmi(srcAmiInfo, targetInfo)
}
