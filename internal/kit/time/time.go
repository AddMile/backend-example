package time

import (
	"fmt"
	"time"
)

const (
	TimeLayout     = "15:04"
	DateLayout     = "2006-01-02"
	DateTimeLayout = "2006-01-02 15:04"
)

var NowUTC = func() time.Time {
	return time.Now().UTC()
}

var NowLocal = func(location *time.Location) time.Time {
	return NowUTC().In(location)
}

func UnixUTC() int64 {
	return NowUTC().Unix()
}

func ConvertLocalDateToUTC(local, timezone string) (time.Time, error) {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, fmt.Errorf("parsing timezone: %w", err)
	}

	parsedTime, err := time.Parse(TimeLayout, local)
	if err != nil {
		return time.Time{}, fmt.Errorf("parsing time: %w", err)
	}

	nowLocal := NowUTC().In(location)
	newLocal := time.Date(nowLocal.Year(), nowLocal.Month(), nowLocal.Day(),
		parsedTime.Hour(), parsedTime.Minute(), 0, 0, location)

	return newLocal.UTC(), nil
}

func ParseDateLocal(date, timezone string) (time.Time, error) {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, err
	}

	parsedTime, err := time.ParseInLocation(DateLayout, date, location)
	if err != nil {
		return time.Time{}, err
	}

	return parsedTime, nil
}

func FormatDateTime(t time.Time) string {
	return t.Format(DateTimeLayout)
}

func IsEqualDate(t1, t2 time.Time) bool {
	layout := "2006-01-02"

	return t1.Format(layout) == t2.Format(layout)
}
