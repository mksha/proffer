package validator

import (
	"fmt"
)

var (
	// TestIsZero
	zeroMap    map[string]int
	zeroStruct struct{ Name string }
	zeroSlice  []int

	// TestCheckRequiredFieldsInStruct

)

// test cases for IsZero function.
var isZeroTestCases = []struct {
	name  string
	value interface{}
	want  bool
}{
	{
		name:  "zero value bool",
		value: false,
		want:  true,
	},
	{
		name:  "non-zero value bool",
		value: false,
		want:  true,
	},
	{
		name:  "zero value int",
		value: 0,
		want:  true,
	},
	{
		name:  "non-zero value int",
		value: 1,
		want:  false,
	},
	{
		name:  "zero value float",
		value: 0.0,
		want:  true,
	},
	{
		name:  "zero value map",
		value: zeroMap,
		want:  true,
	},
	{
		name:  "non-zero value map",
		value: map[int]int{1: 1},
		want:  false,
	},
	{
		name:  "zero value struct",
		value: zeroStruct,
		want:  true,
	},
	{
		name:  "non-zero value struct",
		value: struct{ Name string }{"test"},
		want:  false,
	},
	{
		name:  "zero value slice",
		value: zeroSlice,
		want:  true,
	},
	{
		name:  "non-zero value slice",
		value: []int{1},
		want:  false,
	},
}

// test cases for CheckRequiredFieldsInStruct function.
var checkRequiredFieldsInStructTestCases = []struct {
	name         string
	customStruct interface{}
	index        []int
	want         []error
}{
	{
		name: "struct with required field filled",
		customStruct: struct {
			Name string `required:"true"`
		}{Name: "test"},
		want: make([]error, 0),
	},
	{
		name: "struct with required field filled and having private field",
		customStruct: struct {
			Name string `required:"true"`
			age  int
		}{Name: "test"},
		want: make([]error, 0),
	},
	{
		name: "struct with required field filled, index given and having chain tag",
		customStruct: struct {
			Name string `required:"true" chain:"name"`
		}{Name: "test"},
		index: []int{1},
		want:  make([]error, 0),
	},
	{
		name: "struct with multiple required fields filled",
		customStruct: struct {
			Name string `required:"true"`
			Age  int    `required:"true"`
		}{Name: "test", Age: 1},
		want: make([]error, 0),
	},
	{
		name: "struct with missing required field but having chain tag",
		customStruct: struct {
			Name string `required:"true" chain:"name"`
			Age  int
		}{Age: 1},
		want: []error{fmt.Errorf("name")},
	},
	{
		name: "struct with missing required fields but having chain tag",
		customStruct: struct {
			Name string `required:"true" chain:"name"`
			Age  int    `required:"true" chain:"age"`
		}{},
		want: []error{fmt.Errorf("name"), fmt.Errorf("age")},
	},
	{
		name: "struct with missing required fields and having mapstructure tag instead of chain tag",
		customStruct: struct {
			Name string `required:"true" mapstructure:"name"`
			Age  int
		}{Age: 1},
		want: []error{fmt.Errorf("name")},
	},
	{
		name: "struct with required fields missing and having yaml tag instead of chain tag",
		customStruct: struct {
			Name string `required:"true" yaml:"name"`
			Age  int
		}{Age: 1},
		want: []error{fmt.Errorf("name")},
	},
}

// test cases for IsAWSRoleARN function.
var isAWSRoleARNTestCases = []struct {
	name string
	arn  string
	want bool
}{
	{
		name: "aws role arn with default path",
		arn:  "arn:aws:iam::123456789012:role/test",
		want: true,
	},
	{
		name: "aws role arn with gov partition and default path",
		arn:  "arn:aws-us-gov:iam::123456789012:role/test",
		want: true,
	},
	{
		name: "aws role arn with cn partition and default path",
		arn:  "arn:aws-cn:iam::123456789012:role/test",
		want: true,
	},
	{
		name: "invalid aws role arn with invalid account id",
		arn:  "arn:aws:iam::1232456789012:role/test",
		want: false,
	},
	{
		name: "invalid aws role arn with invalid role path",
		arn:  "arn:aws:iam::123456789012:role/",
		want: false,
	},
	{
		name: "invalid aws role arn with invalid role name",
		arn:  "arn:aws:iam::123456789012:role/&ab%",
		want: false,
	},
	{
		name: "invalid aws role arn with invalid partition",
		arn:  "arn:aws-test:iam::123456789012:role/test",
		want: false,
	},
	{
		name: "invalid aws role arn with more than 64 characters",
		arn:  "arn:aws:iam::123456789012:role/test-hsdfhjgsdhjf-jhhjffsd-3662-234-3432-4324xxxfsdfsdsdgsgdg2deffrg",
		want: false,
	},
}

// test cases for IsAWSRegion function.
var isAWSRegionTestCases = []struct {
	name   string
	region string
	want   bool
}{
	{
		name:   "us-east-2 region",
		region: "us-east-2",
		want:   true,
	},
	{
		name:   "us-east-1 region",
		region: "us-east-1",
		want:   true,
	},
	{
		name:   "us-west-1 region",
		region: "us-west-1",
		want:   true,
	},
	{
		name:   "us-west-2 region",
		region: "us-west-2",
		want:   true,
	},
	{
		name:   "af-south-1 region",
		region: "af-south-1",
		want:   true,
	},
	{
		name:   "ap-northeast-1 region",
		region: "ap-northeast-1",
		want:   true,
	},
	{
		name:   "ca-central-1 region",
		region: "ca-central-1",
		want:   true,
	},
	{
		name:   "cn-north-1 region",
		region: "cn-north-1",
		want:   true,
	},
	{
		name:   "cn-northwest-1 region",
		region: "cn-northwest-1",
		want:   true,
	},
	{
		name:   "eu-central-1 region",
		region: "eu-central-1",
		want:   true,
	},
	{
		name:   "eu-west-1 region",
		region: "eu-west-1",
		want:   true,
	},
	{
		name:   "eu-west-2 region",
		region: "eu-west-2",
		want:   true,
	},
	{
		name:   "eu-south-1 region",
		region: "eu-south-1",
		want:   true,
	},
	{
		name:   "eu-west-3 region",
		region: "eu-west-3",
		want:   true,
	},
	{
		name:   "eu-north-1 region",
		region: "eu-north-1",
		want:   true,
	},
	{
		name:   "me-south-1 region",
		region: "me-south-1",
		want:   true,
	},
	{
		name:   "sa-east-1 region",
		region: "sa-east-1",
		want:   true,
	},
	{
		name:   "us-gov-east-1 region",
		region: "us-gov-east-1",
		want:   true,
	},
	{
		name:   "us-gov-west-1 region",
		region: "us-gov-west-1",
		want:   true,
	},
	{
		name:   "ap-south-1 region",
		region: "ap-south-1",
		want:   true,
	},
	{
		name:   "ap-northeast-3 region",
		region: "ap-northeast-3",
		want:   true,
	},
	{
		name:   "ap-northeast-2 region",
		region: "ap-northeast-2",
		want:   true,
	},
	{
		name:   "ap-southeast-1 region",
		region: "ap-southeast-1",
		want:   true,
	},
	{
		name:   "ap-southeast-2 region",
		region: "ap-southeast-2",
		want:   true,
	},
	{
		name:   "invalid region",
		region: "us-southeast-4",
		want:   false,
	},
	{
		name:   "invalid region",
		region: "us-gov-west-12",
		want:   false,
	},
	{
		name:   "invalid region",
		region: "sus-west-1",
		want:   false,
	},
	{
		name:   "invalid region",
		region: "us-central-1",
		want:   false,
	},
}

// test cases for IsAWSAMIID function.
var isAWSAMIIDTestCases = []struct {
	name string
	id   string
	want bool
}{
	{
		name: "valid ami id with 8 digits",
		id:   "ami-12345678",
		want: true,
	},
	{
		name: "valid ami id with more than 8 digits",
		id:   "ami-1234567891213",
		want: true,
	},
	{
		name: "invalid valid ami id with more than 17 digits",
		id:   "ami-123456789121368365835",
		want: false,
	},
	{
		name: "invalid valid ami id with less than 8 digits",
		id:   "ami-12345",
		want: false,
	},
	{
		name: "invalid valid ami id with no prefix -ami",
		id:   "123456785",
		want: false,
	},
}

// test cases for IsAWSAMIName function.
var isAWSAMINameTestCases = []struct {
	name    string
	amiName string
	want    bool
}{
	{
		name:    "valid and unique ami name",
		amiName: "test-1234-unique-string",
		want:    true,
	},
	{
		name:    "invalid ami name with less than 3 characters",
		amiName: "te",
		want:    false,
	},
	{
		name:    "invalid ami name with invalid characters",
		amiName: "test-&#%$",
		want:    false,
	},
	{
		name:    "invalid ami name with more than 128 characters characters",
		amiName: "test-wufbuwbwiubwiuecbweucbdddeeeeddduwebcuwbcueucbceuwceiuwgeuivgwuivgwuvufhufewfwefwefwefewfwefwefwefewfwefwfwfiwehfiuwhfuwefiuwef",
		want:    false,
	},
}

// test cases for IsAWSTagKey function.
var isAWSTagKeyTestCases = []struct {
	name string
	key  string
	want bool
}{
	{
		name: "valid tag key",
		key:  "Purpose",
		want: true,
	},
	{
		name: "valid tag key with special characters",
		key:  "test - , 8 4 $ *2 @)@*$ JJ f ",
		want: true,
	},
	{
		name: "invalid tag key with no character",
		key:  "",
		want: false,
	},
	{
		name: "invalid tag key with more than 127 characters",
		key:  "h djd test test wg test test test test test test test test test test %$test$ test TEST test TEST TEST ETE TETS TEST TEST TEST test *646466446* ",
		want: false,
	},
}

// test cases for IsAWSTagValue function.
var isAWSTagValueTestCases = []struct {
	name  string
	value string
	want  bool
}{
	{
		name:  "valid tag value",
		value: "test",
		want:  true,
	},
	{
		name:  "valid tag value with zero characters",
		value: "",
		want:  true,
	},
	{
		name:  "valid tag value with special characters",
		value: "test - , 8 4 $ *2 @)@*$ JJ f ",
		want:  true,
	},
	{
		name:  "invalid tag value with more than 128 characters",
		value: "h djd test test wg test test eit eit test test test test test test test %$test$ test TEST test TEST TEST ETE TETS TEST TEST TEST test *646466446* ",
		want:  false,
	},
}

// test cases for IsAWSAccountID function.
var isAWSAccountIDTestCases = []struct {
	name string
	id   string
	want bool
}{
	{
		name: "valid aws account id",
		id:   "123456789012",
		want: true,
	},
	{
		name: "invalid aws account id with less than 12 digits",
		id:   "1234567",
		want: false,
	},
	{
		name: "invalid aws account id with more than 12 digits",
		id:   "123456789102345",
		want: false,
	},
	{
		name: "invalid aws account id with characters",
		id:   "23jj2test",
		want: false,
	},
}
