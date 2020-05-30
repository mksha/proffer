package parser

import (
	"github.com/proffer/components"
)

// test cases for getEnv function.
var getEnvTestCases = []struct {
	name    string
	envName string
	want    string
}{
	{
		name:    "exported env var",
		envName: "TEST_ENV_VAR",
		want:    "TEST_ENV_VAR_VALUE",
	},
	{
		name:    "un-exported env var",
		envName: "TEST_ENV_VAR_NOT_EXIST",
		want:    "",
	},
}

// test cases for setDefaultValue function.
type args struct {
	givenValue   string
	currentValue string
}

var setDefaultValueTestCases = []struct {
	name string
	args args
	want string
}{
	{
		name: "current value is not empty",
		args: args{
			givenValue:   "test2",
			currentValue: "test",
		},
		want: "test",
	},
	{
		name: "current value is empty",
		args: args{
			givenValue:   "test2",
			currentValue: "",
		},
		want: "test2",
	},
}

// test cases for parseTemplate function.
var parseTemplateTestCases = []struct {
	name    string
	dsc     string
	want    string
	wantErr bool
}{
	{
		name:    "valid template config with single resource copyami type",
		dsc:     "../test/valid_config/copyami.yml",
		want:    "output.yml",
		wantErr: false,
	},
	{
		name:    "valid template config with single resource shareami type",
		dsc:     "../test/valid_config/shareami.yml",
		want:    "output.yml",
		wantErr: false,
	},
	{
		name:    "valid template config with multiple resources",
		dsc:     "../test/valid_config/proffer.yml",
		want:    "output.yml",
		wantErr: false,
	},
	{
		name:    "invalid template config with invalid template function",
		dsc:     "../test/invalid_config/sample1.yml",
		want:    "",
		wantErr: true,
	},
	{
		name:    "invalid template config with wrong syntax",
		dsc:     "../test/invalid_config/sample2.yml",
		want:    "",
		wantErr: true,
	},
	{
		name:    "invalid template config with missing template function arguments",
		dsc:     "../test/invalid_config/sample3.yml",
		want:    "",
		wantErr: true,
	},
	{
		name:    "invalid template config with unknown template variable",
		dsc:     "../test/invalid_config/sample4.yml",
		want:    "",
		wantErr: true,
	},
	{
		name:    "invalid template config with invalid nested template",
		dsc:     "../test/invalid_config/sample5.yml",
		want:    "output.yml",
		wantErr: true,
	},
	{
		name:    "template config file does not exist",
		dsc:     "file_does_not_exists",
		want:    "",
		wantErr: true,
	},
}

// test cases for UnmarshalYaml function.
var unmarshalYamlTestCases = []struct {
	name     string
	filePath string
	want     TemplateConfig
	wantErr  bool
}{
	{
		name:     "valid yaml with copyami type resource",
		filePath: "../test/valid_yml/copyami.yml",
		want: TemplateConfig{
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
			},
		},
		wantErr: false,
	},
	{
		name:     "valid yaml config with shareami type resource",
		filePath: "../test/valid_yml/shareami.yml",
		want: TemplateConfig{
			RawResources: []components.RawResource{
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
									"accountId": 121616226324,
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
									"accountId": 121266418583,
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
		name:     "valid yaml config with multiple resource types",
		filePath: "../test/valid_yml/proffer.yml",
		want: TemplateConfig{
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
		name:     "invalid yaml config with missing fields",
		filePath: "../test/invalid_yml/sample1.yml",
		want:     TemplateConfig{},
		wantErr:  true,
	},
	{
		name:     "invalid yaml config with wrong syntax",
		filePath: "../test/invalid_yml/sample2.yml",
		want:     TemplateConfig{},
		wantErr:  true,
	},
	{
		name:     "invalid yaml config with invalid characters",
		filePath: "../test/invalid_yml/sample3.yml",
		want:     TemplateConfig{},
		wantErr:  true,
	},
	{
		name:     "yml file does not exist",
		filePath: "file_path_does_not_exist",
		want:     TemplateConfig{},
		wantErr:  true,
	},
	{
		name:     "yml file with no read permission",
		filePath: "../test/invalid_yml/no_read_permission.yml",
		want:     TemplateConfig{},
		wantErr:  true,
	},
}

// test cases for UnmarshalDefaultVars function.
var unmarshalDefaultVarsTestCases = []struct {
	name            string
	defaultVarsPath string
	wantErr         bool
}{
	{
		name:            "yml file does not exist",
		defaultVarsPath: "file_path_does_not_exist",
		wantErr:         true,
	},
	{
		name:            "invalid yaml config with invalid characters",
		defaultVarsPath: "../test/invalid_yml/sample3.yml",
		wantErr:         false,
	},
	{
		name:            "valid template config with vars defined",
		defaultVarsPath: "../test/valid_config/dynamicproffer.yml",
		wantErr:         false,
	},
	{
		name:            "yml file with no read permission",
		defaultVarsPath: "../test/invalid_yml/no_read_permission.yml",
		wantErr:         true,
	},
}

// test cases for UnmarshalDynamicVars function.
var unmarshalDynamicVarsTestCases = []struct {
	name            string
	dynamicVarsPath string
	wantErr         bool
}{
	{
		name:            "yml file does not exist",
		dynamicVarsPath: "file_path_does_not_exist",
		wantErr:         true,
	},
	{
		name:            "invalid yaml config with invalid characters",
		dynamicVarsPath: "../test/invalid_yml/sample3.yml",
		wantErr:         true,
	},
	{
		name:            "invalid yaml config with invalid characters",
		dynamicVarsPath: "../test/invalid_yml/data2.yml",
		wantErr:         true,
	},
	{
		name:            "valid empty data yaml",
		dynamicVarsPath: "../test/invalid_yml/data1.yml",
		wantErr:         false,
	},
	{
		name:            "valid template config with dynamic defined",
		dynamicVarsPath: "../test/valid_yml/data.yml",
		wantErr:         false,
	},
	{
		name:            "yml file with no read permission",
		dynamicVarsPath: "../test/invalid_yml/no_read_permission.yml",
		wantErr:         true,
	},
}

// test cases for getVars function.
var getVarsTestCases = []struct {
	name    string
	varName string
	want    string
	wantOpt string
}{
	{
		name:    "dynamic var of type string exist in yml config",
		varName: "srcRegion",
		want:    "us-east-1",
	},
	{
		name:    "dynamic var of type list exist in yml config",
		varName: "targetCopyRegions",
		want:    "[us-east-1, us-west-2]",
	},
	{
		name:    "dynamic var of type map exist in yml config",
		varName: "addExtraTags",
		want:    "{CreatedBy: local-testing-tool, Type: testing}",
		wantOpt: "{Type: testing, CreatedBy: local-testing-tool}",
	},
	{
		name:    "dynamic var of type string does not exist in dynamic and default yml config",
		varName: "notexist",
		want:    "<nil>",
	},
	{
		name:    "dynamic var of type string does not exist in dynamic and but exist in default yml config",
		varName: "targetAccountsCommonRegions",
		want:    "[us-east-1, us-west-2]",
	},
}
