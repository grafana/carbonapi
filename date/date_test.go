package date

import (
	"fmt"
	"testing"
	"time"
)

func TestParseAtTimeOr(t *testing.T) {

	defaultTimeZone := time.Local
	// 16 Aug 1994 15:30
	defer MockTimeNow(MockTimeNow(func() time.Time {
		return time.Date(1994, time.August, 16, 15, 30, 0, 100, defaultTimeZone)
	}))

	const shortForm = "15:04 2006-Jan-02"

	var tests = []struct {
		input  string
		output string
	}{
		{"midnight", "00:00 1994-Aug-16"},
		{"noon", "12:00 1994-Aug-16"},
		{"teatime", "16:00 1994-Aug-16"},
		{"tomorrow", "00:00 1994-Aug-17"},

		{"noon 08/12/94", "12:00 1994-Aug-12"},
		{"midnight 20060812", "00:00 2006-Aug-12"},
		{"noon tomorrow", "12:00 1994-Aug-17"},

		{"17:04 19940812", "17:04 1994-Aug-12"},
		{"-1day", "15:30 1994-Aug-15"},
		{"19940812", "00:00 1994-Aug-12"},

		{"today-2d", "00:00 1994-Aug-14"},
		{"today-1h", "23:00 1994-Aug-15"},
		{"yesterday+12h", "12:00 1994-Aug-15"},
		{"now-1h", "14:30 1994-Aug-16"},
		{"now+30min", "16:00 1994-Aug-16"},
		{"noon+3h", "15:00 1994-Aug-16"},
		{"midnight-30min", "23:30 1994-Aug-15"},

		// case-insensitive
		{"NOW", "15:30 1994-Aug-16"},
		{"Today-1h", "23:00 1994-Aug-15"},
		{"MIDNIGHT", "00:00 1994-Aug-16"},

		// 4-digit year in MM/DD/YYYY
		{"01/02/2014", "00:00 2014-Jan-02"},
		{"noon 08/12/2006", "12:00 2006-Aug-12"},
	}

	for _, tt := range tests {
		got := ParseAtTimeOr(tt.input, "Local", defaultTimeZone, 0)
		ts, err := time.ParseInLocation(shortForm, tt.output, defaultTimeZone)
		if err != nil {
			panic(fmt.Sprintf("error parsing time: %q: %v", tt.output, err))
		}

		want := int64(ts.Unix())
		if got != want {
			t.Errorf("ParseAtTimeOr(%q)=%v, want %v", tt.input, got, want)
		}
	}
}
