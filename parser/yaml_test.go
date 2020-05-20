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
