package parser

import (
	"log"
	"os"
	"testing"
)

func Test_getEnv(t *testing.T) {
	err := os.Setenv("TEST_ENV_VAR", "TEST_ENV_VAR_VALUE")
	if err != nil {
		log.Fatal(err)
	}

	for n := range getEnvTestCases {
		tt := getEnvTestCases[n]
		t.Run(tt.name, func(t *testing.T) {
			if got := getEnv(tt.envName); got != tt.want {
				t.Errorf("getEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Benchmark_getEnv(b *testing.B) {
	// bench combined time to run through all test cases
	for n := range getEnvTestCases {
		tc := getEnvTestCases[n]

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// ignoring errors and results because we're just timing function execution
				_ = getEnv(tc.envName)
			}
		})
	}
}

func Test_setDefaultValue(t *testing.T) {
	for n := range setDefaultValueTestCases {
		tt := setDefaultValueTestCases[n]
		t.Run(tt.name, func(t *testing.T) {
			if got := setDefaultValue(tt.args.givenValue, tt.args.currentValue); got != tt.want {
				t.Errorf("setDefaultValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Benchmark_setDefaultValue(b *testing.B) {
	// bench combined time to run through all test cases
	for n := range setDefaultValueTestCases {
		tc := setDefaultValueTestCases[n]

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// ignoring errors and results because we're just timing function execution
				_ = setDefaultValue(tc.args.givenValue, tc.args.currentValue)
			}
		})
	}
}

func TestParseTemplate(t *testing.T) {
	for n := range parseTemplateTestCases {
		tt := parseTemplateTestCases[n]
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTemplate(tt.dsc)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkParseTemplate(b *testing.B) {
	// bench combined time to run through all test cases
	for n := range parseTemplateTestCases {
		tc := parseTemplateTestCases[n]

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// ignoring errors and results because we're just timing function execution
				_, _ = ParseTemplate(tc.dsc)
			}
		})
	}
}

func Test_getVar(t *testing.T) {
	if err := UnmarshalDynamicVars("../test/valid_yml/data.yml"); err != nil {
		t.Errorf("got err during parsing of dynamic vars: \n%v", err)
	}

	if err := UnmarshalDefaultVars("../test/valid_config/dynamicproffer.yml"); err != nil {
		t.Errorf("got err during parsing of default vars: \n%v", err)
	}

	for n := range getVarsTestCases {
		tt := getVarsTestCases[n]
		t.Run(tt.name, func(t *testing.T) {
			if got := getVar(tt.varName); got != tt.want {
				t.Errorf("getVar() = %v, want %v", got, tt.want)
			}
		})
	}
}
