package cmd

import (
	"github.com/proffer/components"
	"github.com/proffer/parser"
)

var tests = []struct {
	name    string
	dsc     string
	want    parser.TemplateConfig
	wantErr bool
}{
	{
		name: "valid config with multiple resource types",
		dsc:  "../test/valid_yml/proffer.yml",
		want: parser.TemplateConfig{
			RawResources: []components.RawResource{
				{
					Name: "Copy AMI To Commercial Account Regions",
					Type: "aws-copyami",
					Config: map[string]interface{}{
						"source": map[interface{}]interface{}{
							"amiFilters": map[interface{}]interface{}{
								"name":           "test-image",
								"tag:Department": "DevOps",
								"tag:Purpose":    "Testing",
							},
							"profile": "demo-2",
							"region":  "us-east-1",
						},
						"target": map[interface{}]interface{}{
							"addExtraTags": map[interface{}]interface{}{
								"CreatedBy": "local-testing-tool",
								"Type":      "testing",
							},
							"copyTagsAcrossRegions": true,
							"regions": []interface{}{
								"ap-northeast-1",
								"ap-northeast-2",
								"us-west-2",
							},
						},
					},
				},
				{
					Name: "Share AMI With Other Accounts",
					Type: "aws-shareami",
					Config: map[string]interface{}{
						"source": map[interface{}]interface{}{
							"amiFilters": map[interface{}]interface{}{
								"name":           "test-image",
								"tag:Department": "DevOps",
								"tag:Purpose":    "Testing",
							},
							"profile": "demo-2",
						},
						"target": map[interface{}]interface{}{
							"accountRegionMappingList": []interface{}{
								map[interface{}]interface{}{
									"accountId": 591616226324,
									"profile":   "demo-1",
									"regions": []interface{}{
										"ap-northeast-1",
									},
									"copyTagsAcrossAccounts": true,
									"addExtraTags": map[interface{}]interface{}{
										"CreatedBy": "SharedByDemo1",
										"Type":      "AMITesting",
										"Home":      "/home/test",
									},
								},
								map[interface{}]interface{}{
									"accountId": 221266418583,
									"profile":   "platform-aws",
									"regions": []interface{}{
										"ap-northeast-2",
									},
								},
							},
							"copyTagsAcrossAccounts":    true,
							"addCreateVolumePermission": true,
							"commonRegions": []interface{}{
								"us-east-1",
								"us-west-2",
							},
						},
					},
				},
			},
		},
		wantErr: false,
	},
	{
		name:    "invalid config with missing fields",
		dsc:     "../test/invalid_yml/sample1.yml",
		want:    parser.TemplateConfig{},
		wantErr: true,
	},
	{
		name:    "non existing template config",
		dsc:     "not_exist.yml",
		want:    parser.TemplateConfig{},
		wantErr: true,
	},
}

// test cases for validCmd command.
var validaCmdTestCases = []struct {
	name    string
	args    []string
	want    string
	wantErr bool
}{
	{
		name:    "valid subcommand and template",
		args:    []string{"validate", "../test/valid_yml/proffer.yml"},
		want:    "",
		wantErr: false,
	},
	{
		name:    "valid subcommand with template and debug argument",
		args:    []string{"validate", "../test/valid_yml/proffer.yml", "-d"},
		want:    "",
		wantErr: false,
	},
	{
		name:    "missing template file",
		args:    []string{"validate"},
		want:    "proffer template file is missing: Pls pass proffer template file to apply",
		wantErr: true,
	},
	{
		name:    "missing template file with config flag",
		args:    []string{"validate", "--config", "../test/valid_yml/proffer.yml"},
		want:    "proffer template file is missing: Pls pass proffer template file to apply",
		wantErr: true,
	},
	{
		name:    "invalid template file",
		args:    []string{"validate", "../test/invalid_config/sample1.yml"},
		want:    "line 5: field ss not found in type components.RawResource",
		wantErr: true,
	},
	{
		name:    "zero resources",
		args:    []string{"validate", "../test/invalid_config/sample2.yml"},
		want:    "empty resource found in list 'resources'",
		wantErr: true,
	},
	{
		name:    "config with missing keys",
		args:    []string{"validate", "not_exist.yml"},
		want:    "no such file or directory",
		wantErr: true,
	},
	{
		name:    "config with missing resource name",
		args:    []string{"validate", "../test/invalid_config/sample3.yml"},
		want:    "Error: name",
		wantErr: true,
	},
	{
		name:    "config with invalid resource type",
		args:    []string{"validate", "../test/invalid_config/sample4.yml"},
		want:    "invalid resource type",
		wantErr: true,
	},
}

var applyCmdTestCases = []struct {
	name    string
	args    []string
	want    string
	wantErr bool
}{
	// {
	// 	name:    "apply without creds",
	// 	args:    []string{"apply", "../test/valid_yml/proffer.yml"},
	// 	want:    "",
	// 	wantErr: true,
	// },
	{
		name:    "config with invalid resource type",
		args:    []string{"apply", "../test/invalid_config/sample4.yml"},
		want:    "invalid resource type",
		wantErr: true,
	},
	{
		name:    "invalid template file",
		args:    []string{"apply", "../test/invalid_config/sample1.yml"},
		want:    "line 5: field ss not found in type components.RawResource",
		wantErr: true,
	},
	{
		name:    "missing template file",
		args:    []string{"apply"},
		want:    "proffer Configuration file missing",
		wantErr: true,
	},
	{
		name:    "config with missing keys",
		args:    []string{"apply", "not_exist.yml"},
		want:    "no such file or directory",
		wantErr: true,
	},
}
