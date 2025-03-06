package money_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AddMile/backend/internal/kit/money"
)

func TestToMicros(t *testing.T) {
	type testCase struct {
		name  string
		value string
		want  int64
	}

	testCases := []testCase{
		{
			name:  "converts 19.99 to 19,990,000",
			value: "19.99",
			want:  19_990_000,
		},
		{
			name:  "converts 14.99 to 14,990,000",
			value: "14.99",
			want:  14_990_000,
		},
		{
			name:  "converts 20.00 to 20,000,000",
			value: "20.00",
			want:  20_000_000,
		},
		{
			name:  "converts 00.00 to 0",
			value: "00.00",
			want:  0,
		},
		{
			name:  "converts -20.00 to -20,000,000",
			value: "-20.00",
			want:  -20_000_000,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := money.ToMicros(tc.value)
			assert.NoError(t, err)

			assert.Equal(t, tc.want, result, tc.value)
		})
	}
}
