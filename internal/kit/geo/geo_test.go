package geo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLookupCountryCode(t *testing.T) {
	db, err := New()
	assert.NoError(t, err)

	tests := []struct {
		name     string
		ip       string
		expected string
	}{
		{
			name:     "Valid IP from Ukraine",
			ip:       "178.165.42.67",
			expected: "UA",
		},
		{
			name:     "Invalid IP address",
			ip:       "169.254.169.126",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := db.LookupCountryCode(tt.ip)
			assert.Equal(t, tt.expected, result)
		})
	}
}
