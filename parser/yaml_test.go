package parser

import (
	"reflect"
	"testing"
)

func TestUnmarshalYaml(t *testing.T) {
	for n := range unmarshalYamlTestCases {
		tt := unmarshalYamlTestCases[n]
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalYaml(tt.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalYaml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnmarshalYaml() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkUnmarshalYaml(b *testing.B) {
	// bench combined time to run through all test cases
	for n := range unmarshalYamlTestCases {
		tc := unmarshalYamlTestCases[n]

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// ignoring errors and results because we're just timing function execution
				_, _ = UnmarshalYaml(tc.filePath)
			}
		})
	}
}
