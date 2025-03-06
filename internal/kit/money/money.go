package money

import (
	"strconv"
)

// microExp is the multiplier used for converting between floating-point values
// and microseconds. Its value represents one million, as there are 1,000,000 microseconds in a second.
const microExp = 1e6

// ToMicros converts a string value representing seconds into a Micros value
// (time duration in microseconds).
func ToMicros(v string) (int64, error) {
	parsed, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0, err
	}

	return int64(parsed * microExp), nil
}
