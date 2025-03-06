package time

import (
	"testing"
	"time"
)

func WithMockedNowUTC(t *testing.T, hours int) {
	t.Helper()

	stubTime := time.Date(2024, 7, 21, hours, 0, 0, 0, time.UTC)
	originalNow := NowUTC
	NowUTC = func() time.Time { return stubTime }

	t.Cleanup(func() {
		NowUTC = originalNow
	})
}
