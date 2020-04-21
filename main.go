package main

import (
	"os"
)

func main() {
	cmdArgs := os.Args
	dsc := cmdArgs[1]
	execute(dsc)
	// srcAmiInfo := SrcAmiInfo{
	// 	Region: aws.String("us-west-2"),
	// 	Filters: []*ec2.Filter{
	// 		{
	// 			Name: aws.String("image-id"),
	// 			Values: []*string{
	// 				aws.String("ami-01b55974057720a43"),
	// 			},
	// 		},
	// 	},
	// }

	// targetInfo := TargetInfo{
	// 	Regions: []*string{
	// 		aws.String("ap-southeast-1"),
	// 	},
	// }
	// copyAmi(srcAmiInfo, targetInfo)
}
