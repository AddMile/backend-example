package time_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	timekit "github.com/AddMile/backend/internal/kit/time"
)

func TestConvertLocalDateToUTC(t *testing.T) {
	timekit.WithMockedNowUTC(t, 12)

	type testCase struct {
		name     string
		timezone string
		local    string
		expected time.Time
	}

	testCases := []testCase{
		{
			name:     "given local time and timezone, should return datetime in UTC",
			timezone: "Europe/Kyiv",
			local:    "12:00",
			expected: time.Date(2024, 7, 21, 9, 0, 0, 0, time.UTC),
		},
		{
			name:     "should return datetime followed the day before",
			timezone: "Europe/Kyiv",
			local:    "02:00",
			expected: time.Date(2024, 7, 20, 23, 0, 0, 0, time.UTC),
		},
		{
			name:     "should return datetime followed the next day",
			timezone: "America/Sao_Paulo", // -3 GMT
			local:    "22:00",
			expected: time.Date(2024, 7, 22, 01, 0, 0, 0, time.UTC),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := timekit.ConvertLocalDateToUTC(tc.local, tc.timezone)
			assert.NoError(t, err)

			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestIsEqualDate(t *testing.T) {
	tests := []struct {
		name     string
		t1       time.Time
		t2       time.Time
		expected bool
	}{
		{
			name:     "Same date, different times",
			t1:       time.Date(2025, 2, 26, 10, 30, 0, 0, time.UTC),
			t2:       time.Date(2025, 2, 26, 23, 59, 59, 999999999, time.UTC),
			expected: true,
		},
		{
			name:     "Different dates",
			t1:       time.Date(2025, 2, 26, 10, 30, 0, 0, time.UTC),
			t2:       time.Date(2025, 2, 27, 0, 0, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "Same date, different timezones",
			t1:       time.Date(2025, 2, 26, 10, 0, 0, 0, time.UTC),
			t2:       time.Date(2025, 2, 26, 10, 0, 0, 0, time.FixedZone("Custom", 3600)),
			expected: true,
		},
		{
			name:     "Edge case - Midnight",
			t1:       time.Date(2025, 2, 26, 0, 0, 0, 0, time.UTC),
			t2:       time.Date(2025, 2, 26, 0, 0, 0, 1, time.UTC),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := timekit.IsEqualDate(tt.t1, tt.t2)
			assert.Equal(t, tt.expected, result)
		})
	}
}
