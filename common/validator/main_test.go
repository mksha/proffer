package validator

import (
	"reflect"
	"testing"
)

func TestIsZero(t *testing.T) {
	for n := range isZeroTestCases {
		tt := isZeroTestCases[n]
		t.Run(tt.name, func(t *testing.T) {
			if got := IsZero(tt.value); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkIsZero(b *testing.B) {
	// bench combined time to run through all test cases
	for n := range isZeroTestCases {
		tc := isZeroTestCases[n]

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// ignoring errors and results because we're just timing function execution
				_ = IsZero(tc.value)
			}
		})
	}
}

func TestCheckRequiredFieldsInStruct(t *testing.T) {
	for n := range checkRequiredFieldsInStructTestCases {
		tt := checkRequiredFieldsInStructTestCases[n]
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckRequiredFieldsInStruct(tt.customStruct, tt.index...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CheckRequiredFieldsInStruct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkCheckRequiredFieldsInStruct(b *testing.B) {
	// bench combined time to run through all test cases
	for n := range checkRequiredFieldsInStructTestCases {
		tc := checkRequiredFieldsInStructTestCases[n]

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// ignoring errors and results because we're just timing function execution
				_ = CheckRequiredFieldsInStruct(tc.customStruct, tc.index...)
			}
		})
	}
}

func TestIsAWSRoleARN(t *testing.T) {
	for n := range isAWSRoleARNTestCases {
		tt := isAWSRoleARNTestCases[n]
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAWSRoleARN(tt.arn); got != tt.want {
				t.Errorf("IsAWSRoleARN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkIsAWSRoleARN(b *testing.B) {
	// bench combined time to run through all test cases
	for n := range isAWSRoleARNTestCases {
		tc := isAWSRoleARNTestCases[n]

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// ignoring errors and results because we're just timing function execution
				_ = IsAWSRoleARN(tc.arn)
			}
		})
	}
}

func TestIsAWSRegion(t *testing.T) {
	for n := range isAWSRegionTestCases {
		tt := isAWSRegionTestCases[n]
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAWSRegion(tt.region); got != tt.want {
				t.Errorf("IsAWSRegion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkIsAWSRegion(b *testing.B) {
	// bench combined time to run through all test cases
	for n := range isAWSRegionTestCases {
		tc := isAWSRegionTestCases[n]

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// ignoring errors and results because we're just timing function execution
				_ = IsAWSRegion(tc.region)
			}
		})
	}
}

func TestIsAWSAMIID(t *testing.T) {
	for n := range isAWSAMIIDTestCases {
		tt := isAWSAMIIDTestCases[n]
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAWSAMIID(tt.id); got != tt.want {
				t.Errorf("IsAWSAMIID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkIsAWSAMIID(b *testing.B) {
	// bench combined time to run through all test cases
	for n := range isAWSAMIIDTestCases {
		tc := isAWSAMIIDTestCases[n]

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// ignoring errors and results because we're just timing function execution
				_ = IsAWSAMIID(tc.id)
			}
		})
	}
}

func TestIsAWSAMIName(t *testing.T) {
	for n := range isAWSAMINameTestCases {
		tt := isAWSAMINameTestCases[n]
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAWSAMIName(tt.amiName); got != tt.want {
				t.Errorf("IsAWSAMIName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkIsAWSAMIName(b *testing.B) {
	// bench combined time to run through all test cases
	for n := range isAWSAMINameTestCases {
		tc := isAWSAMINameTestCases[n]

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// ignoring errors and results because we're just timing function execution
				_ = IsAWSAMIName(tc.amiName)
			}
		})
	}
}

func TestIsAWSTagKey(t *testing.T) {
	for n := range isAWSTagKeyTestCases {
		tt := isAWSTagKeyTestCases[n]
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAWSTagKey(tt.key); got != tt.want {
				t.Errorf("IsAWSTagKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkIsAWSTagKey(b *testing.B) {
	// bench combined time to run through all test cases
	for n := range isAWSTagKeyTestCases {
		tc := isAWSTagKeyTestCases[n]

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// ignoring errors and results because we're just timing function execution
				_ = IsAWSTagKey(tc.key)
			}
		})
	}
}

func TestIsAWSTagValue(t *testing.T) {
	for n := range isAWSTagValueTestCases {
		tt := isAWSTagValueTestCases[n]
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAWSTagValue(tt.value); got != tt.want {
				t.Errorf("IsAWSTagValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkIsAWSTagValue(b *testing.B) {
	// bench combined time to run through all test cases
	for n := range isAWSTagValueTestCases {
		tc := isAWSTagValueTestCases[n]

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// ignoring errors and results because we're just timing function execution
				_ = IsAWSTagValue(tc.value)
			}
		})
	}
}

func TestIsAWSAccountID(t *testing.T) {
	for n := range isAWSAccountIDTestCases {
		tt := isAWSAccountIDTestCases[n]
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAWSAccountID(tt.id); got != tt.want {
				t.Errorf("IsAWSAccountID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkIsAWSAccountID(b *testing.B) {
	// bench combined time to run through all test cases
	for n := range isAWSAccountIDTestCases {
		tc := isAWSAccountIDTestCases[n]

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// ignoring errors and results because we're just timing function execution
				_ = IsAWSAccountID(tc.id)
			}
		})
	}
}
